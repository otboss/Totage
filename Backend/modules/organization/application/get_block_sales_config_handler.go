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

func (s *Handlers) GetBlockSalesConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetBlockSalesConfigQuery],
) (*connect.Response[endpoints.GetBlockSalesConfigResponse], error) {
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

	blockConfig, err := db_organization.Query.GetBlockSalesConfig(ctx, req.Msg.OrganizationId)
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
	if blockConfig.MinOrderQuantity.Valid != false {
		minOrderQuantity = &blockConfig.MinOrderQuantity.Float64
	}

	return connect.NewResponse(&endpoints.GetBlockSalesConfigResponse{
		BlockPriceMin:    blockConfig.BlockPriceMin,
		BlockPricePerKm:  blockConfig.BlockPricePerKm,
		BlockPrice:       blockConfig.BlockPrice,
		AvailableSupply:  blockConfig.AvailableSupply,
		MinOrderQuantity: minOrderQuantity,
		IsEnabled:        blockConfig.IsEnabled.Bool,
	}), nil
}
