package application

import (
	"context"
	"errors"
	"fmt"
	"time"
	"waterordersystem/GOBACKEND/modules/organization/domain/templates"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	db_shared "waterordersystem/GOBACKEND/shared/infrastructure"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) ForgotPasswordHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ForgotPasswordCommand],
) (*connect.Response[endpoints.ForgotPasswordResponse], error) {

	if utils.Env.CAPTCHA_DISABLED == false {
		recaptchaValidation, err := clients.GoogleCloudApi.RecaptchaValidation(clients.RecaptchaEventRequest{
			Event: clients.RecaptchaEvent{
				Token:          req.Msg.RecaptchaEvent.Token,
				ExpectedAction: req.Msg.RecaptchaEvent.ExpectedAction,
				SiteKey:        req.Msg.RecaptchaEvent.SiteKey,
			},
		})
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		if recaptchaValidation.TokenProperties.Valid == false {
			return nil, errors.New("recaptcha validation failed")
		}
	}

	user, err := db_user.Query.GetUserByEmail(ctx, req.Msg.Email)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	verificationSession, err := db_user.Query.GetVerificationSessionByUserId(ctx, user.ID)
	if err == pgx.ErrNoRows || (err == nil && time.Since(verificationSession.CreatedAt.Time) >= 30*time.Minute) {
		verificationSession, err := db_shared.Query.CreateVerificationSession(ctx)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		_, err = db_user.Query.CreateVerificationSession(ctx, db_user.CreateVerificationSessionParams{
			UserID:              user.ID,
			VerificationCodesID: verificationSession.ID,
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
				[]string{user.Email.String},
				"Password Reset Requested",
				templates.ForgotPassword(verificationSession.ID.String(), verificationSession.VerificationCode),
			)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeInternal,
					errors.New("internal server error"),
				)
			}
		}

		return connect.NewResponse(&endpoints.ForgotPasswordResponse{
			SessionToken: verificationSession.ID.String(),
		}), nil
	}

	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return nil, errors.New("verification session already in progress")
}
