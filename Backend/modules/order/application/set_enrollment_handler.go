package application

import (
	"context"
	"errors"
	"slices"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) SetEnrollmentHandler(
	ctx context.Context,
	req *connect.Request[endpoints.SetEnrollmentCommand],
) (*connect.Response[endpoints.SetEnrollmentResponse], error) {
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

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: vehicle.OrganizationID,
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

	_, err = db_order.Query.UpdateEnrollment(ctx, db_order.UpdateEnrollmentParams{
		VehicleID:       req.Msg.VehicleId,
		Freight:         slices.Contains(req.Msg.OrderTypes, common.OrderType_FREIGHT),
		WaterPotable:    slices.Contains(req.Msg.OrderTypes, common.OrderType_WATER_PURIFIED),
		WaterUnpurified: slices.Contains(req.Msg.OrderTypes, common.OrderType_WATER_UNPURIFIED),
		Ice:             slices.Contains(req.Msg.OrderTypes, common.OrderType_ICE),
		Blocks:          slices.Contains(req.Msg.OrderTypes, common.OrderType_BLOCK),
		Coconut:         slices.Contains(req.Msg.OrderTypes, common.OrderType_COCONUT),
		Egg:             slices.Contains(req.Msg.OrderTypes, common.OrderType_EGG),
		Charcoal:        slices.Contains(req.Msg.OrderTypes, common.OrderType_CHARCOAL),
		CookingGas:      slices.Contains(req.Msg.OrderTypes, common.OrderType_COOKING_GAS),
		Honey:           slices.Contains(req.Msg.OrderTypes, common.OrderType_HONEY),
		Cesspool:        slices.Contains(req.Msg.OrderTypes, common.OrderType_CESSPOOL),
		Construction:    slices.Contains(req.Msg.OrderTypes, common.OrderType_CONSTRUCTION),
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.SetEnrollmentResponse{}), nil
}
