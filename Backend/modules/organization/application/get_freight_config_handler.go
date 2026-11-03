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

func (s *Handlers) GetFreightConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetFreightConfigQuery],
) (*connect.Response[endpoints.GetFreightConfigResponse], error) {
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

	freightConfig, err := db_organization.Query.GetFreightConfig(ctx, req.Msg.OrganizationId)
	if err == pgx.ErrNoRows {
		return connect.NewResponse(&endpoints.GetFreightConfigResponse{}), nil
	}
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var sLat *float64 = nil
	var sLon *float64 = nil
	var lLat *float64 = nil
	var lLon *float64 = nil

	if freightConfig.SolidWasteDumpSiteLatitude != -999.0 && freightConfig.SolidWasteDumpSiteLongitude != -999.0 {
		sLat = &freightConfig.SolidWasteDumpSiteLatitude
		sLon = &freightConfig.SolidWasteDumpSiteLongitude
	}

	if freightConfig.LiquidWasteDumpSiteLatitude != -999.0 && freightConfig.LiquidWasteDumpSiteLongitude != -999.0 {
		lLat = &freightConfig.LiquidWasteDumpSiteLatitude
		lLon = &freightConfig.LiquidWasteDumpSiteLongitude
	}

	return connect.NewResponse(&endpoints.GetFreightConfigResponse{
		SolidWasteDumpSiteLatitude:   sLat,
		SolidWasteDumpSiteLongitude:  sLon,
		LiquidWasteDumpSiteLatitude:  lLat,
		LiquidWasteDumpSiteLongitude: lLon,
		IsEnabled:                    freightConfig.IsEnabled.Bool,
	}), nil
}
