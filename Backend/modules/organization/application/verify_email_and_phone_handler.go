package application

import (
	"context"
	"errors"
	"time"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	db_shared "waterordersystem/GOBACKEND/shared/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) VerifyEmailAndPhoneHandler(
	ctx context.Context,
	req *connect.Request[endpoints.VerifyEmailAndPhoneCommand],
) (*connect.Response[endpoints.VerifyEmailAndPhoneResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: req.Msg.OrganizationId,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if account.IsAdmin == false {
		return nil, errors.New("insufficient account permission")
	}

	organizationVerificationCode, err := db_organization.Query.GetOrganizationVerificationCodeByOrganizationId(ctx, req.Msg.OrganizationId)
	if err == pgx.ErrNoRows || (err == nil && time.Since(organizationVerificationCode.CreatedAt.Time) >= 30*time.Minute) {
		emailVerificationSession, err := db_shared.Query.CreateVerificationSession(ctx)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		phoneVerificationSession, err := db_shared.Query.CreateVerificationSession(ctx)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		_, err = db_organization.Query.CreateOrganizationVerificationCode(ctx, db_organization.CreateOrganizationVerificationCodeParams{
			OrganizationID:           req.Msg.OrganizationId,
			EmailVerificationCodesID: emailVerificationSession.ID,
			PhoneVerificationCodesID: phoneVerificationSession.ID,
		})

		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		return connect.NewResponse(&endpoints.VerifyEmailAndPhoneResponse{
			IsVerificationComplete:        false,
			EmailVerificationSessionToken: emailVerificationSession.ID.String(),
			PhoneVerificationSessionToken: phoneVerificationSession.ID.String(),
		}), nil
	}

	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if req.Msg.VerificationSession == nil {
		return nil, errors.New("verification session param required for verification")
	}

	if req.Msg.VerificationSession.EmailVerificationSession.SessionToken != organizationVerificationCode.EmailVerificationCodesID.String() {
		return nil, errors.New("invalid email verification session token provided")
	}

	if req.Msg.VerificationSession.PhoneVerificationSession.SessionToken != organizationVerificationCode.PhoneVerificationCodesID.String() {
		return nil, errors.New("invalid phone verification session token provided")
	}

	u, err := uuid.Parse(req.Msg.VerificationSession.EmailVerificationSession.SessionToken)
	if err != nil {
		return nil, errors.New("bad email session token provided")
	}

	isEmailVerificationSessionSuccessful, err := db_shared.Query.IsVerificationSessionSuccessful(ctx, db_shared.IsVerificationSessionSuccessfulParams{
		ID: pgtype.UUID{
			Bytes: u,
			Valid: true,
		},
		VerificationCode: req.Msg.VerificationSession.EmailVerificationSession.VerificationCode,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if isEmailVerificationSessionSuccessful == false {
		return nil, errors.New("failed verification attempt or too many failed attempts occurred")
	}

	u, err = uuid.Parse(req.Msg.VerificationSession.PhoneVerificationSession.SessionToken)
	if err != nil {
		return nil, errors.New("bad phone session token provided")
	}

	isPhoneVerificationSessionSuccessful, err := db_shared.Query.IsVerificationSessionSuccessful(ctx, db_shared.IsVerificationSessionSuccessfulParams{
		ID: pgtype.UUID{
			Bytes: u,
			Valid: true,
		},
		VerificationCode: req.Msg.VerificationSession.PhoneVerificationSession.VerificationCode,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if isPhoneVerificationSessionSuccessful == false {
		return nil, errors.New("failed verification attempt or too many failed attempts occurred")
	}

	err = db_organization.Query.UpdateOrganizationVerifyEmailAndPhone(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, err = db_organization.Query.CreateOrganizationMailList(ctx, db_organization.CreateOrganizationMailListParams{
		OrganizationID: req.Msg.OrganizationId,
		Email:          account.Email.String,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.VerifyEmailAndPhoneResponse{
		IsVerificationComplete: true,
	}), nil
}
