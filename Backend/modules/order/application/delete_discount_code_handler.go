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

func (s *Handlers) DeleteDiscountCodeHandler(
	ctx context.Context,
	req *connect.Request[endpoints.DeleteDiscountCodeCommand],
) (*connect.Response[endpoints.DeleteDiscountCodeResponse], error) {

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

	err = db_order.Query.DeleteDiscountCode(ctx, req.Msg.DiscountCodeId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.DeleteDiscountCodeResponse{}), nil
}
