package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) ApproveVehicleHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ApproveVehicleCommand],
) (*connect.Response[endpoints.ApproveVehicleResponse], error) {

	// TODO: Check if current user is authorized to view vehicles

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	_, err = db_user.Query.GetSystemAdministratorByUserId(ctx, claims.UserID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	// TODO: Check if user is superuser

	// TODO: If user super user approve vehicle

	// TODO: Log vehicle approval

	return connect.NewResponse(&endpoints.ApproveVehicleResponse{}), nil
}
