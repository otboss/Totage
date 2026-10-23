package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"waterordersystem/GOBACKEND/modules/organization/domain/templates"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	db_shared "waterordersystem/GOBACKEND/shared/infrastructure"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) CreateUserViaEmailHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateUserViaEmailCommand],
) (*connect.Response[endpoints.CreateUserViaEmailResponse], error) {

	emailVerificationCode, err := db_user.Query.GetEmailVerificationCodeByEmail(ctx, req.Msg.Email)
	if err == pgx.ErrNoRows || (err == nil && time.Since(emailVerificationCode.CreatedAt.Time) >= 30*time.Minute) {
		newVerificationSession, err := db_shared.Query.CreateVerificationSession(ctx)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		_, err = db_user.Query.CreateEmailVerificationCode(ctx, db_user.CreateEmailVerificationCodeParams{
			Email:               req.Msg.Email,
			VerificationCodesID: newVerificationSession.ID,
		})
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		if utils.Env.EMAIL_DISABLED == false {
			err = clients.EmailApi.SendEmail(
				fmt.Sprintf("\"%s\" <%s>", utils.Env.SMTP_NAME, utils.Env.SMTP_USER),
				[]string{req.Msg.Email},
				"Account Verification",
				templates.EmailAddressVerificationCode(newVerificationSession.ID.String(), newVerificationSession.VerificationCode),
			)
			if err != nil {
				log.Println(err.Error())
			}
		}

		return connect.NewResponse(&endpoints.CreateUserViaEmailResponse{
			VerificationSessionToken: newVerificationSession.ID.String(),
			JwtToken:                 nil,
		}), nil
	}

	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if req.Msg.UserInfo == nil {
		return nil, errors.New("user info required for verification")
	}

	if req.Msg.UserInfo.VerificationSession.SessionToken != emailVerificationCode.VerificationCodesID.String() {
		return nil, errors.New("invalid verification session token provided")
	}

	u, err := uuid.Parse(req.Msg.UserInfo.VerificationSession.SessionToken)
	if err != nil {
		return nil, errors.New("bad session token provided")
	}

	isVerificationSessionSuccessful, err := db_shared.Query.IsVerificationSessionSuccessful(ctx, db_shared.IsVerificationSessionSuccessfulParams{
		ID: pgtype.UUID{
			Bytes: u,
			Valid: true,
		},
		VerificationCode: req.Msg.UserInfo.VerificationSession.VerificationCode,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if isVerificationSessionSuccessful == false {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("failed verification attempt or too many failed attempts occurred"))
	}

	id, err := db_user.Query.CreateUser(ctx, db_user.CreateUserParams{
		FirstName: req.Msg.UserInfo.FirstName,
		LastName:  req.Msg.UserInfo.LastName,
		Phone:     req.Msg.UserInfo.Phone,
		Email: pgtype.Text{
			String: req.Msg.Email,
			Valid:  true,
		},
		LastIpAddress: req.Peer().Addr,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code != "23505" {
				// not duplicate key error
				return nil, connect.NewError(
					connect.CodeInternal,
					errors.New("internal server error"),
				)
			}
			usr, err := db_user.Query.GetUserByEmail(ctx, req.Msg.Email)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInternal,
					errors.New("internal server error"),
				)
			}
			id = usr.ID
		}
	}

	jwtToken, err := utils.JWT.CreateToken(id, nil, utils.Env.ENC_KEY)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	user, err := db_user.Query.GetUserByEmail(ctx, req.Msg.Email)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	passwordResetverificationSession, err := db_shared.Query.CreateVerificationSession(ctx)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, err = db_user.Query.CreateVerificationSession(ctx, db_user.CreateVerificationSessionParams{
		UserID:              user.ID,
		VerificationCodesID: passwordResetverificationSession.ID,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.CreateUserViaEmailResponse{
		VerificationSessionToken: passwordResetverificationSession.ID.String(),
		JwtToken:                 &jwtToken,
		VerificationCode:         &passwordResetverificationSession.VerificationCode,
	}), nil
}
