package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) UpdateTruckLocationHandler(
	ctx context.Context,
	req *connect.Request[endpoints.UpdateTruckLocationCommand],
) (*connect.Response[endpoints.UpdateTruckLocationResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	organization, err := db_organization.Query.GetOrganizationByTruckId(ctx, req.Msg.TruckId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: organization.ID,
	})
	if err != nil {
		return nil, errors.New("account not found")
	}

	if account.IsAdmin == false && account.IsDriver == false {
		return nil, errors.New("insufficient account permission")
	}

	err = db_logistic.Query.UpdateTruckLocation(ctx, db_logistic.UpdateTruckLocationParams{
		ID:        req.Msg.TruckId,
		Longitude: req.Msg.Coordinates.Longitude,
		Latitude:  req.Msg.Coordinates.Latitude,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.UpdateTruckLocationResponse{}), nil
}
