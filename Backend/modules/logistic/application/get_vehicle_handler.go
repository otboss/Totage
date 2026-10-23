package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) GetVehicleHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetVehicleQuery],
) (*connect.Response[endpoints.GetVehicleResponse], error) {

	// TODO: Check if current user is authorized to view vehicles

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

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

	_, err = db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: organization.ID,
	})
	if err == pgx.ErrNoRows {
		return nil, errors.New("account not found")
	}
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var truckId *int64 = nil
	var rentalVehicleId *int64 = nil

	rentalVehicle, err := db_logistic.Query.GetRentalVehicleByVehicleId(ctx, vehicle.ID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}
	if err != pgx.ErrNoRows {
		rentalVehicleId = &rentalVehicle.ID
	}

	truck, err := db_logistic.Query.GetTruckByVehicleId(ctx, vehicle.ID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}
	if err != pgx.ErrNoRows {
		truckId = &truck.ID
	}

	return connect.NewResponse(&endpoints.GetVehicleResponse{
		Vehicle: &endpoints.Vehicle{
			VehicleId:            vehicle.ID,
			Year:                 int64(vehicle.Year),
			Make:                 vehicle.Make,
			Model:                vehicle.Model,
			ChassisNumber:        vehicle.ChassisNumber,
			LicensePlateNumber:   vehicle.LicensePlateNumber,
			SeatingCapacity:      int64(vehicle.SeatingCapacity),
			FuelType:             endpoints.FuelType(endpoints.FuelType_value[string(vehicle.FuelType)]),
			IsManualTransmission: vehicle.IsManualTransmission,
			RentalVehicleId:      rentalVehicleId,
			TruckId:              truckId,
			DriverAccountId:      vehicle.DriverAccountID.String(),
			IsApproved:           vehicle.IsApproved.Bool,
		},
	}), nil
}
