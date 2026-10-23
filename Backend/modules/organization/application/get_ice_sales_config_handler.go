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

func (s *Handlers) GetIceSalesConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetIceSalesConfigQuery],
) (*connect.Response[endpoints.GetIceSalesConfigResponse], error) {
	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	_, err = db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: req.Msg.OrganizationId,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("not permitted")
		}
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	iceConfig, err := db_organization.Query.GetIceSalesConfig(ctx, req.Msg.OrganizationId)
	if err == pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodeNotFound,
			errors.New("config not found"),
		)
	}
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var minOrderQuantity *float64 = nil
	if iceConfig.MinOrderQuantity.Valid != false {
		minOrderQuantity = &iceConfig.MinOrderQuantity.Float64
	}

	return connect.NewResponse(&endpoints.GetIceSalesConfigResponse{
		IcePriceMin:      iceConfig.IcePriceMin,
		IcePricePerKm:    iceConfig.IcePricePerKm,
		IcePrice:         iceConfig.IcePrice,
		AvailableSupply:  iceConfig.AvailableSupply,
		MinOrderQuantity: minOrderQuantity,
		IsEnabled:        iceConfig.IsEnabled.Bool,
	}), nil
}
