package application

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"

	"connectrpc.com/connect"
)

func (s *Handlers) GetRentalVehicleOrderTotalHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetRentalVehicleOrderTotalQuery],
) (*connect.Response[endpoints.GetOrderTotalResponse], error) {

	rentalVehicle, err := db_logistic.Query.GetRentalVehicle(ctx, req.Msg.RentalVehicleId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	total, err := HelperFunctions.CalculateRentalVehicleOrderTotal(
		db_logistic.LogisticRentalVehicle{
			MinDaysRental:   rentalVehicle.MinDaysRental,
			DailyRentalCost: rentalVehicle.DailyRentalCost,
		},
		time.UnixMilli(req.Msg.StartDate),
		time.UnixMilli(req.Msg.EndDate),
	)
	if err != nil {
		return nil, err
	}

	if req.Msg.DiscountCode != nil {
		total, err = HelperFunctions.GetDiscountValue(ctx, total, *req.Msg.DiscountCode)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
	}

	organization, err := db_organization.Query.GetOrganization(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	response, err := HelperFunctions.BuildOrderTotalResponse(
		ctx,
		db_organization.SharedOrderTypeVEHICLERENTAL,
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

func (h *helperFunctions) CalculateRentalVehicleOrderTotal(rentalVehicle db_logistic.LogisticRentalVehicle, startDate time.Time, endDate time.Time) (float64, error) {
	rentalDuration := startDate.Sub(endDate)
	rentalDifferenceInDays := math.Abs(rentalDuration.Hours() / 24)

	if int(rentalDifferenceInDays) < int(rentalVehicle.MinDaysRental) {
		return 0, errors.New(fmt.Sprintf("min rental duration in days is %d", rentalVehicle.MinDaysRental))
	}
	return rentalDifferenceInDays * rentalVehicle.DailyRentalCost, nil
}
