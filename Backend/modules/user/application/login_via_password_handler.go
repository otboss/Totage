package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Handlers) LoginViaPasswordHandler(
	ctx context.Context,
	req *connect.Request[endpoints.LoginViaPasswordCommand],
) (*connect.Response[endpoints.LoginViaPasswordResponse], error) {

	user, err := db_user.Query.GetUserByEmail(ctx, req.Msg.Email)
	if err == pgx.ErrNoRows {
		return nil, errors.New("invalid email or password")
	}

	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Msg.Password)) != nil {
		// TODO: implement rate limiting to prevent brute forceing attacks
		return nil, errors.New("invalid username or password")
	}

	jwtToken, err := utils.JWT.CreateToken(user.ID, nil, utils.Env.ENC_KEY)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.LoginViaPasswordResponse{
		JwtToken: jwtToken,
	}), nil
}
