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

func (s *Handlers) GetPotableWaterSalesConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetPotableWaterSalesConfigQuery],
) (*connect.Response[endpoints.GetPotableWaterSalesConfigResponse], error) {
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

	potableWaterConfig, err := db_organization.Query.GetPotableWaterSalesConfig(ctx, req.Msg.OrganizationId)
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

	var bottleCost *float64 = nil
	if potableWaterConfig.BottleCost.Valid != false {
		bottleCost = &potableWaterConfig.BottleCost.Float64
	}

	var minOrderQuantity *float64 = nil
	if potableWaterConfig.MinOrderQuantity.Valid != false {
		minOrderQuantity = &potableWaterConfig.MinOrderQuantity.Float64
	}

	return connect.NewResponse(&endpoints.GetPotableWaterSalesConfigResponse{
		WaterPriceMin:    potableWaterConfig.WaterPriceMin,
		WaterPricePerKm:  potableWaterConfig.WaterPricePerKm,
		WaterPrice:       potableWaterConfig.WaterPrice,
		BottleCost:       bottleCost,
		AvailableSupply:  potableWaterConfig.AvailableSupply,
		MinOrderQuantity: minOrderQuantity,
		IsEnabled:        potableWaterConfig.IsEnabled.Bool,
	}), nil
}
