package application

import (
	"context"
	"errors"
	"time"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	db_shared "waterordersystem/GOBACKEND/shared/infrastructure"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (s *Handlers) ResetPasswordHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ResetPasswordCommand],
) (*connect.Response[endpoints.ResetPasswordResponse], error) {

	u, err := uuid.Parse(string(req.Msg.SessionToken))
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	verificationSession, err := db_user.Query.GetVerificationSessionBySessionToken(ctx, pgtype.UUID{
		Bytes: u,
		Valid: true,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if time.Since(verificationSession.CreatedAt.Time) >= 5*time.Minute {
		return nil, errors.New("verification session expired")
	}

	u, err = uuid.Parse(req.Msg.SessionToken)
	if err != nil {
		return nil, errors.New("bad session token provided")
	}

	isVerificationSessionSuccessful, err := db_shared.Query.IsVerificationSessionSuccessful(ctx, db_shared.IsVerificationSessionSuccessfulParams{
		ID: pgtype.UUID{
			Bytes: u,
			Valid: true,
		},
		VerificationCode: req.Msg.VerificationCode,
	})

	if isVerificationSessionSuccessful == false {
		return nil, errors.New("verification attempt failed")
	}

	hPassword, err := bcrypt.GenerateFromPassword([]byte(req.Msg.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	hashedPassword := pgtype.Text{
		String: string(hPassword),
		Valid:  true,
	}

	err = db_user.Query.UpdateUserPassword(ctx, db_user.UpdateUserPasswordParams{
		ID:           verificationSession.UserID,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.ResetPasswordResponse{}), nil
}
