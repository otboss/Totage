package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	template "waterordersystem/GOBACKEND/modules/organization/domain/templates"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func (s *Handlers) CreateAccountHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateAccountCommand],
) (*connect.Response[endpoints.CreateAccountResponse], error) {
	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		OrganizationID: req.Msg.OrganizationId,
		UserID:         claims.UserID,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if account.IsAdmin == false {
		return nil, errors.New("not permitted")
	}

	user, err := db_user.Query.GetUserByEmail(ctx, req.Msg.Email)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	accountId, err := db_organization.Query.CreateAccount(ctx, db_organization.CreateAccountParams{
		UserID:            user.ID,
		OrganizationID:    req.Msg.OrganizationId,
		IsAdmin:           req.Msg.Roles.IsAdmin,
		IsCustomerSupport: req.Msg.Roles.IsDriver,
		IsDriver:          req.Msg.Roles.IsCustomerSupport,
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
			fmt.Sprintf("Organization Invitation - "),
			template.OrganizationInvitation(accountId.String()),
		)
		if err != nil {
			log.Println(err.Error())
		}
	}

	return connect.NewResponse(&endpoints.CreateAccountResponse{
		AccountId: accountId.String(),
	}), nil
}
