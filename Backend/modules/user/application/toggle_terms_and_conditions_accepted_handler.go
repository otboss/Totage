package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"

	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) ToggleTermsAndConditionsAcceptedHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ToggleTermsAndConditionsAcceptedCommand],
) (*connect.Response[endpoints.ToggleTermsAndConditionsAcceptedResponse], error) {
	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	err = db_user.Query.ToggleTermsAndConditionsAccepted(ctx, claims.UserID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.ToggleTermsAndConditionsAcceptedResponse{}), nil
}
