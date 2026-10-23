package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) GetAccountsHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetAccountsQuery],
) (*connect.Response[endpoints.GetAccountsResponse], error) {
	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	_, err = db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: req.Msg.OrganizationId,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("not permitted")
		}
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	accounts, err := db_organization.Query.GetAccounts(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var response endpoints.GetAccountsResponse

	for _, account := range accounts {
		response.Accounts = append(response.Accounts, &endpoints.Account{
			AccountId:         account.ID.String(),
			UserId:            account.UserID,
			FirstName:         account.FirstName,
			LastName:          account.LastName,
			Phone:             account.Phone,
			Email:             account.Email.String,
			IsAdmin:           account.IsAdmin,
			IsCustomerSupport: account.IsCustomerSupport,
			IsDriver:          account.IsDriver,
		})
	}

	return connect.NewResponse(&response), nil
}
