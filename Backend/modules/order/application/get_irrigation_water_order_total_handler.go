package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) GetIrrigationWaterOrderTotalHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetIrrigationWaterOrderTotalQuery],
) (*connect.Response[endpoints.GetOrderTotalResponse], error) {

	organization, err := db_organization.Query.GetOrganization(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if organization.IsApproved == false {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	err = HelperFunctions.CheckOrganizationSubscription(ctx, organization)
	if err != nil {
		return nil, err
	}

	waterConfig, err := db_organization.Query.GetIrrigationWaterSalesConfig(ctx, organization.ID)
	if err == pgx.ErrNoRows || (err == nil && waterConfig.IsEnabled.Bool == false) {
		return nil, errors.New("order type is not enabled for organization")
	}
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	destinationLocation := clients.Location{
		Address:     req.Msg.DeliveryAddress,
		Coordinates: nil,
	}

	distanceCalculation, err := clients.GoogleCloudApi.DistanceBetweenLocations(
		clients.Location{
			Address:     organization.Address,
			Coordinates: nil,
		},
		destinationLocation,
	)

	if distanceCalculation.Rows[0].Elements[0].Status == "ZERO_RESULTS" || distanceCalculation.Rows[0].Elements[0].Status == "NOT_FOUND" {
		return nil, errors.New("an invalid address was provided")
	}

	jobRoundTripDistance := float64(distanceCalculation.Rows[0].Elements[0].Distance.Value/1000) * 2

	total := HelperFunctions.CalculateIrrigationWaterOrderTotal(waterConfig, organization.DutyFreeKm, jobRoundTripDistance, req.Msg.QuantityOfBarrels, req.Msg.QuantityOfWater)

	response, err := HelperFunctions.BuildOrderTotalResponse(
		ctx,
		db_organization.SharedOrderTypeWATERUNPURIFIED, // TODO: REview this type
		organization,
		total,
		req.Msg.DiscountCode,
	)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(response), nil
}

func (h *helperFunctions) CalculateIrrigationWaterOrderTotal(waterConfig db_organization.OrganizationIrrigationWaterSalesConfig, dutyFreeKm float64, jobRoundTripDistance float64, quantityOfBarrels int32, quantityOfWater int32) float64 {
	total := float64(
		(quantityOfBarrels*int32(waterConfig.BarrelCost.Float64) +
			// quantityOfPurifiedWater*int32(waterConfig.PurifiedWaterPricePerGallon.Float64) +
			quantityOfWater*int32(waterConfig.WaterPrice) +
			quantityOfBarrels*int32(waterConfig.BarrelCost.Float64)),
	)
	if jobRoundTripDistance > dutyFreeKm {
		total += jobRoundTripDistance*waterConfig.WaterPricePerKm + waterConfig.WaterPriceMin
	}

	return total
}
