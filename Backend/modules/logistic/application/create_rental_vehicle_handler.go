package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) CreateRentalVehicleHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateRentalVehicleCommand],
) (*connect.Response[endpoints.CreateRentalVehicleResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, _ := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)

	vehicle, err := db_logistic.Query.GetVehicle(ctx, req.Msg.VehicleId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	organization, err := db_organization.Query.GetOrganization(ctx, vehicle.OrganizationID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: organization.ID,
	})
	if err != nil {
		return nil, errors.New("account not found")
	}

	if account.IsAdmin == false {
		return nil, errors.New("insufficient account permission")
	}

	exchangeRateUsd, err := clients.CurrencyApi.GetExchangeRateUsd(string(organization.Currency))
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var driverCost *float64
	if req.Msg.RentalVehicle.DriverCost != nil {
		driverCost = req.Msg.RentalVehicle.DriverCost
	}

	rentalVehicleId, err := db_logistic.Query.CreateRentalVehicle(ctx, db_logistic.CreateRentalVehicleParams{
		VehicleID: vehicle.ID,
		IsDriverMandatory: pgtype.Bool{
			Bool:  req.Msg.RentalVehicle.IsDriverMandatory,
			Valid: true,
		},
		DailyRentalCost: float64(req.Msg.RentalVehicle.DailyRentalCost) / float64(exchangeRateUsd),
		DriverCost: pgtype.Float8{
			Float64: *driverCost,
			Valid:   true,
		},
		MinDaysRental: int16(req.Msg.RentalVehicle.MinDaysRental),
		IsPublic:      req.Msg.RentalVehicle.IsPublic,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	// TODO: Look into why this is needed
	// _, err = db_logistic.Query.CreateVehicleOrderTypeAllocation(ctx, db_logistic.CreateVehicleOrderTypeAllocationParams{
	// 	VehicleID:      req.Msg.VehicleId,
	// 	OrganizationID: req.Msg.OrganizationId,
	// 	OrderType:      db_logistic.SharedOrderTypeVEHICLERENTAL,
	// 	MaxCapacity:    float64(req.Msg.SeatingCapacity),
	// })
	// if err != nil {
	// 	return nil, connect.NewError(
	// 		connect.CodeInternal,
	// 		errors.New("internal server error"),
	// 	)
	// }

	return connect.NewResponse(&endpoints.CreateRentalVehicleResponse{
		RentalVehicleId: rentalVehicleId,
	}), nil
}
