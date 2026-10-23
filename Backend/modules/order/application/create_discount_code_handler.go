package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) CreateDiscountCodeHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateDiscountCodeCommand],
) (*connect.Response[endpoints.CreateDiscountCodeResponse], error) {

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
		return nil, errors.New("insufficient account permission")
	}

	discountCodeType, exists := endpoints.DiscountCodeType_name[int32(req.Msg.DiscountCodeType)]
	if exists == false {
		return nil, errors.New("invalid discount code type provided")
	}

	discountCodeId, err := db_order.Query.CreateDiscountCode(ctx, db_order.CreateDiscountCodeParams{
		Code:             req.Msg.DiscountCode,
		DiscountCodeType: db_order.OrderDisconutCodeType(discountCodeType),
		Value:            req.Msg.Value,
	})

	return connect.NewResponse(&endpoints.CreateDiscountCodeResponse{
		Id: discountCodeId,
	}), nil
}
