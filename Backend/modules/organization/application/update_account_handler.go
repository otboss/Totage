package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) UpdateAccountHandler(
	ctx context.Context,
	req *connect.Request[endpoints.UpdateAccountCommand],
) (*connect.Response[endpoints.UpdateAccountResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	u, err := uuid.Parse(req.Msg.AccountId)
	if err != nil {
		return nil, errors.New("bad account id provided")
	}

	accountId := pgtype.UUID{
		Bytes: u,
		Valid: true,
	}

	account, err := db_organization.Query.GetAccount(ctx, accountId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	currentUserAccount, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: account.OrganizationID,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if currentUserAccount.IsAdmin == false {
		return nil, errors.New("not permitted")
	}

	err = db_organization.Query.UpdateAccountRoles(ctx, db_organization.UpdateAccountRolesParams{
		AccountID:         accountId,
		IsAdmin:           req.Msg.Roles.IsAdmin,
		IsDriver:          req.Msg.Roles.IsDriver,
		IsCustomerSupport: req.Msg.Roles.IsCustomerSupport,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.UpdateAccountResponse{}), nil
}
