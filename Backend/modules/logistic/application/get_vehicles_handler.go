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

func (s *Handlers) GetVehiclesHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetVehiclesQuery],
) (*connect.Response[endpoints.GetVehiclesResponse], error) {

	// TODO: Check if current user is authorized to view vehicles

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	organization, err := db_organization.Query.GetOrganization(ctx, req.Msg.OrganizationId)
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

	var offsetId int64 = 0
	if req.Msg.Pagination.OffsetId != nil {
		offsetId = *req.Msg.Pagination.OffsetId
	}

	filterText := ""
	if req.Msg.Pagination.FilterText != nil {
		filterText = *req.Msg.Pagination.FilterText
	}

	vehicles, err := db_logistic.Query.GetVehicles(ctx, db_logistic.GetVehiclesParams{
		OrganizationID: req.Msg.OrganizationId,
		FilterText:     filterText,
		OffsetID:       offsetId,
		Limitation:     int32(req.Msg.Pagination.Limit),
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var vehiclesResponse []*endpoints.Vehicle

	for _, vehicle := range vehicles {
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

		vehiclesResponse = append(vehiclesResponse, &endpoints.Vehicle{
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
		})
	}

	return connect.NewResponse(&endpoints.GetVehiclesResponse{
		Vehicles: vehiclesResponse,
	}), nil
}
