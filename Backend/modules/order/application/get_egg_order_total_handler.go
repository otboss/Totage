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

func (s *Handlers) GetEggOrderTotalHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetEggOrderTotalQuery],
) (*connect.Response[endpoints.GetOrderTotalResponse], error) {

	organization, err := db_organization.Query.GetOrganization(ctx, int64(req.Msg.OrganizationId))
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

	eggConfig, err := db_organization.Query.GetEggSalesConfig(ctx, organization.ID)
	if err == pgx.ErrNoRows || (err == nil && eggConfig.IsEnabled.Bool == false) {
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

	total := HelperFunctions.CalculateEggOrderTotal(eggConfig, organization.DutyFreeKm, jobRoundTripDistance, req.Msg.QuantityOfEggs)

	response, err := HelperFunctions.BuildOrderTotalResponse(
		ctx,
		db_organization.SharedOrderTypeEGG,
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

func (h *helperFunctions) CalculateEggOrderTotal(eggConfig db_organization.OrganizationEggSalesConfig, dutyFreeKm float64, jobRoundTripDistance float64, quantityOfEggs int32) float64 {
	total := float64(quantityOfEggs)*eggConfig.EggPrice + eggConfig.EggPriceMin
	if jobRoundTripDistance > dutyFreeKm {
		total += jobRoundTripDistance * eggConfig.EggPriceMin
	}
	return total
}
