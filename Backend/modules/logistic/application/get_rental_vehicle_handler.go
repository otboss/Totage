package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) GetRentalVehicleHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetRentalVehicleQuery],
) (*connect.Response[endpoints.GetRentalVehicleResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	rentalVehicle, err := db_logistic.Query.GetRentalVehicleByVehicleId(ctx, req.Msg.VehicleId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	vehicle, err := db_logistic.Query.GetVehicle(ctx, rentalVehicle.VehicleID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	// vehicleOrderTypeAllocation, err := db_logistic.Query.GetVehicleOrderTypeAllocation(ctx, db_logistic.GetVehicleOrderTypeAllocationParams{
	// 	VehicleID:      rentalVehicle.VehicleID,
	// 	OrganizationID: rentalVehicle.OrganizationID,
	// 	OrderType:      db_logistic.SharedOrderTypeVEHICLERENTAL,
	// })
	// if err != nil {
	// 	return nil, connect.NewError(
	// 		connect.CodeInternal,
	// 		errors.New("internal server error"),
	// 	)
	// }

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

	return connect.NewResponse(&endpoints.GetRentalVehicleResponse{
		RentalVehicle: &endpoints.RentalVehicle{
			DailyRentalCost:   rentalVehicle.DailyRentalCost,
			MinDaysRental:     int64(rentalVehicle.MinDaysRental),
			DriverCost:        &rentalVehicle.DriverCost.Float64,
			IsDriverMandatory: rentalVehicle.IsDriverMandatory.Bool,
		},
	}), nil
}
