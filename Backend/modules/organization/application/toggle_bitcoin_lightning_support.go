package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) ToggleBitcoinLightningSupportHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ToggleBitcoinLightningSupportCommand],
) (*connect.Response[endpoints.ToggleBitcoinLightningSupportResponse], error) {
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

	bitcoinLightningNotToggledOn := true
	_, err = db_organization.Query.GetBitcoinLightningCredentialByOrganizationID(ctx, req.Msg.OrganizationId)
	if err == pgx.ErrNoRows {
		bitcoinLightningNotToggledOn = false

		lightningAddressEncrypted, err := utils.AES.Encrypt(utils.Env.ENC_KEY, req.Msg.LightningAddress)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		_, err = db_organization.Query.CreateBitcoinLightningCredential(ctx, db_organization.CreateBitcoinLightningCredentialParams{
			OrganizationID:            req.Msg.OrganizationId,
			LightningAddressEncrypted: lightningAddressEncrypted,
		})
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
	} else if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if bitcoinLightningNotToggledOn {
		err = db_organization.Query.DeleteBitcoinLightningCredentialByOrganizationID(ctx, req.Msg.OrganizationId)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
	}

	return connect.NewResponse(&endpoints.ToggleBitcoinLightningSupportResponse{}), nil
}
