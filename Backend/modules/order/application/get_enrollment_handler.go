package application

import (
	"context"
	"errors"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) GetEnrollmentHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetEnrollmentQuery],
) (*connect.Response[endpoints.GetEnrollmentResponse], error) {
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

	enrollments, err := db_order.Query.GetEnrollmentByVehicleId(ctx, req.Msg.VehicleId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var orderTypes []common.OrderType

	if enrollments.Blocks {
		orderTypes = append(orderTypes, common.OrderType_BLOCK)
	}
	if enrollments.Freight {
		orderTypes = append(orderTypes, common.OrderType_FREIGHT)
	}
	if enrollments.Charcoal {
		orderTypes = append(orderTypes, common.OrderType_CHARCOAL)
	}
	if enrollments.Coconut {
		orderTypes = append(orderTypes, common.OrderType_COCONUT)
	}
	if enrollments.CookingGas {
		orderTypes = append(orderTypes, common.OrderType_COOKING_GAS)
	}
	if enrollments.Egg {
		orderTypes = append(orderTypes, common.OrderType_EGG)
	}
	if enrollments.Cesspool {
		orderTypes = append(orderTypes, common.OrderType_CESSPOOL)
	}
	if enrollments.Construction {
		orderTypes = append(orderTypes, common.OrderType_CONSTRUCTION)
	}
	if enrollments.Honey {
		orderTypes = append(orderTypes, common.OrderType_HONEY)
	}
	if enrollments.Ice {
		orderTypes = append(orderTypes, common.OrderType_ICE)
	}
	if enrollments.WaterPotable {
		orderTypes = append(orderTypes, common.OrderType_WATER_PURIFIED)
	}
	if enrollments.WaterUnpurified {
		orderTypes = append(orderTypes, common.OrderType_WATER_UNPURIFIED)
	}

	return connect.NewResponse(&endpoints.GetEnrollmentResponse{
		OrderTypes: orderTypes,
	}), nil
}
