package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) GetUserInfoHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetUserInfoQuery],
) (*connect.Response[endpoints.GetUserInfoResponse], error) {
	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid authorization token provided"),
		)
	}

	user, err := db_user.Query.GetUser(ctx, claims.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, connect.NewError(
				connect.CodeNotFound,
				errors.New("user not found"),
			)
		}
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var userEmail *string
	if user.Email.Valid != false {
		userEmail = &user.Email.String
	}

	return connect.NewResponse(&endpoints.GetUserInfoResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     userEmail,
		Phone:     user.Phone,
	}), nil
}
