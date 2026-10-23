package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) CreateVehicleHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateVehicleCommand],
) (*connect.Response[endpoints.CreateVehicleResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, _ := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)

	organization, err := db_organization.Query.GetOrganization(ctx, req.Msg.OrganizationId)
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

	driverAccBytes, err := uuid.Parse(req.Msg.DriverAccountId)
	if err != nil {
		return nil, errors.New("bad session token provided")
	}

	driverAccountId := pgtype.UUID{
		Bytes: driverAccBytes,
		Valid: true,
	}

	vehicleId, err := db_logistic.Query.CreateVehicle(ctx, db_logistic.CreateVehicleParams{
		OrganizationID:       req.Msg.OrganizationId,
		Year:                 int16(req.Msg.Year),
		Make:                 req.Msg.Make,
		Model:                req.Msg.Model,
		ChassisNumber:        req.Msg.ChassisNumber,
		DriverAccountID:      driverAccountId,
		LicensePlateNumber:   req.Msg.LicensePlateNumber,
		SeatingCapacity:      int16(req.Msg.SeatingCapacity),
		FuelType:             db_logistic.LogisticVehicleFuelType(endpoints.FuelType_name[int32(req.Msg.FuelType)]),
		IsManualTransmission: req.Msg.IsManualTransmission,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	// TODO: Do vehicle enrollment validation
	// TODO: Add vehicle enrollment in _order.enrollment table

	return connect.NewResponse(&endpoints.CreateVehicleResponse{
		VehicleId: vehicleId,
	}), nil
}
