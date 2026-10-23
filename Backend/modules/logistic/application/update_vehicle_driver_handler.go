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

func (s *Handlers) UpdateVehicleDriverHandler(
	ctx context.Context,
	req *connect.Request[endpoints.UpdateVehicleDriverCommand],
) (*connect.Response[endpoints.UpdateVehicleDriverResponse], error) {

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

	driverAccBytes, err := uuid.Parse(req.Msg.DriverAccountId)
	if err != nil {
		return nil, errors.New("bad session token provided")
	}

	driverAccountId := pgtype.UUID{
		Bytes: driverAccBytes,
		Valid: true,
	}

	err = db_logistic.Query.UpdateVehicleDriver(ctx, db_logistic.UpdateVehicleDriverParams{
		VehicleID:       req.Msg.VehicleId,
		DriverAccountID: driverAccountId,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.UpdateVehicleDriverResponse{}), nil
}
