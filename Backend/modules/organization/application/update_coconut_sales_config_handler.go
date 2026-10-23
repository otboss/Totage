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

func (s *Handlers) UpdateCoconutSalesConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.UpdateCoconutSalesConfigCommand],
) (*connect.Response[endpoints.UpdateCoconutSalesConfigResponse], error) {
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

	err = db_organization.Query.UpdateOrganizationCoconutSalesConfig(ctx, db_organization.UpdateOrganizationCoconutSalesConfigParams{
		OrganizationID:    req.Msg.OrganizationId,
		CoconutPriceMin:   req.Msg.CoconutPriceMin,
		CoconutPrice:      req.Msg.CoconutPrice,
		CoconutPricePerKm: req.Msg.CoconutPricePerKm,
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
	return connect.NewResponse(&endpoints.UpdateCoconutSalesConfigResponse{}), nil
}
