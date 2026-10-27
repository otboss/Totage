package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) GetEnabledPaymentMethodsHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetEnabledPaymentMethodsQuery],
) (*connect.Response[endpoints.GetEnabledPaymentMethodsResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	_, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	var enabledPaymentMethods []common.PaymentMethod

	_, err = db_organization.Query.GetBitcoinCredentialByOrganizationID(ctx, req.Msg.OrganizationId)
	if err == pgx.ErrNoRows {
		enabledPaymentMethods = append(enabledPaymentMethods, common.PaymentMethod_BITCOIN)
	} else if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, err = db_organization.Query.GetPayPalCredentialByOrganizationID(ctx, req.Msg.OrganizationId)
	if err == pgx.ErrNoRows {
		enabledPaymentMethods = append(enabledPaymentMethods, common.PaymentMethod_PAYPAL)
	} else if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, err = db_organization.Query.GetBitcoinLightningCredentialByOrganizationID(ctx, req.Msg.OrganizationId)
	if err == pgx.ErrNoRows {
		enabledPaymentMethods = append(enabledPaymentMethods, common.PaymentMethod_BITCOIN_LIGHTNING)
	} else if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.GetEnabledPaymentMethodsResponse{
		PaymentMethods: enabledPaymentMethods,
	}), nil
}
