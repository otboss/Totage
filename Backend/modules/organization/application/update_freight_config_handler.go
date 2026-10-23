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

func (s *Handlers) UpdateFreightConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.UpdateFreightConfigCommand],
) (*connect.Response[endpoints.UpdateFreightConfigResponse], error) {
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

	err = db_organization.Query.UpdateOrganizationFreightConfig(ctx, db_organization.UpdateOrganizationFreightConfigParams{
		OrganizationID: req.Msg.OrganizationId,
		IsEnabled: pgtype.Bool{
			Bool:  req.Msg.IsEnabled,
			Valid: true,
		},
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.UpdateFreightConfigResponse{}), nil
}
