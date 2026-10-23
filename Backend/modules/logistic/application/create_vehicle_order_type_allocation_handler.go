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

func (s *Handlers) CreateVehicleOrderTypeAllocationHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateVehicleOrderTypeAllocationCommand],
) (*connect.Response[endpoints.CreateVehicleOrderTypeAllocationResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: req.Msg.OrganizationId,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if account.IsAdmin == false {
		return nil, connect.NewError(
			connect.CodePermissionDenied,
			errors.New("not permitted"),
		)
	}

	id, err := db_logistic.Query.CreateVehicleOrderTypeAllocation(ctx, db_logistic.CreateVehicleOrderTypeAllocationParams{
		VehicleID:      req.Msg.VehicleId,
		OrganizationID: req.Msg.OrganizationId,
		OrderType:      db_logistic.SharedOrderType(req.Msg.OrderType.String()),
		MaxCapacity:    req.Msg.MaxCapacity,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.CreateVehicleOrderTypeAllocationResponse{
		VehicleAllocationId: id,
	}), nil
}
