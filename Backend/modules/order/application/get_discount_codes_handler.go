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

func (s *Handlers) GetDiscountCodesHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetDiscountCodesQuery],
) (*connect.Response[endpoints.GetDiscountCodesResponse], error) {

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

	var id int64 = 0
	if req.Msg.Pagination.OffsetId != nil {
		id = *req.Msg.Pagination.OffsetId
	}

	discountCodes, err := db_order.Query.GetAllDiscountCodes(ctx, db_order.GetAllDiscountCodesParams{
		OffsetID:       id,
		OrganizationID: req.Msg.OrganizationId,
		Limit:          int32(req.Msg.Pagination.Limit),
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var response []*endpoints.GetDiscountCodesResponse_DiscountCodes
	for _, discountCode := range discountCodes {

		discountCodeType, exists := endpoints.DiscountCodeType_value[string(discountCode.DiscountCodeType)]
		if exists == false {
			// TODO: Log error
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		response = append(response, &endpoints.GetDiscountCodesResponse_DiscountCodes{
			DiscountCodeType: endpoints.DiscountCodeType(discountCodeType),
			Value:            discountCode.Value,
			DiscountCode:     discountCode.Code,
			CreatedAt:        discountCode.CreatedAt.Time.String(),
		})
	}

	return connect.NewResponse(&endpoints.GetDiscountCodesResponse{
		DiscountCodes: response,
	}), nil
}
