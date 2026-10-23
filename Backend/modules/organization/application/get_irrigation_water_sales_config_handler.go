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

func (s *Handlers) GetIrrigationWaterSalesConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetIrrigationWaterSalesConfigQuery],
) (*connect.Response[endpoints.GetIrrigationWaterSalesConfigResponse], error) {
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

	irrigationWaterConfig, err := db_organization.Query.GetIrrigationWaterSalesConfig(ctx, req.Msg.OrganizationId)
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

	var barrelCost *float64 = nil
	if irrigationWaterConfig.BarrelCost.Valid != false {
		barrelCost = &irrigationWaterConfig.BarrelCost.Float64
	}

	var minOrderQuantity *float64 = nil
	if irrigationWaterConfig.MinOrderQuantity.Valid != false {
		minOrderQuantity = &irrigationWaterConfig.MinOrderQuantity.Float64
	}

	return connect.NewResponse(&endpoints.GetIrrigationWaterSalesConfigResponse{
		WaterPriceMin:    irrigationWaterConfig.WaterPriceMin,
		WaterPricePerKm:  irrigationWaterConfig.WaterPricePerKm,
		WaterPrice:       irrigationWaterConfig.WaterPrice,
		BarrelCost:       barrelCost,
		AvailableSupply:  irrigationWaterConfig.AvailableSupply,
		MinOrderQuantity: minOrderQuantity,
		IsEnabled:        irrigationWaterConfig.IsEnabled.Bool,
	}), nil
}
