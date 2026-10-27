package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) TogglePaypalSupportHandler(
	ctx context.Context,
	req *connect.Request[endpoints.TogglePaypalSupportCommand],
) (*connect.Response[endpoints.TogglePaypalSupportResponse], error) {
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
		return nil, errors.New("not permitted")
	}

	paypalRestSdkClientIdEncrypted, err := utils.AES.Encrypt(utils.Env.ENC_KEY, req.Msg.PaypalRestSdkClientId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	paypalRestSdkSecretEncrypted, err := utils.AES.Encrypt(utils.Env.ENC_KEY, req.Msg.PaypalRestSdkSecret)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, err = db_organization.Query.UpsertPayPalCredential(ctx, db_organization.UpsertPayPalCredentialParams{
		OrganizationID: req.Msg.OrganizationId,
		Email:          "",
		RestSdkClientIDEncrypted: pgtype.Text{
			String: paypalRestSdkClientIdEncrypted,
			Valid:  true,
		},
		RestSdkSecretEncrypted: pgtype.Text{
			String: paypalRestSdkSecretEncrypted,
			Valid:  true,
		},
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.TogglePaypalSupportResponse{}), nil
}
