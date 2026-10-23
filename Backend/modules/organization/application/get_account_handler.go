package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) GetAccountHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetAccountQuery],
) (*connect.Response[endpoints.GetAccountResponse], error) {

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

	u, err := uuid.Parse(req.Msg.AccountId)
	if err != nil {
		return nil, errors.New("bad account id provided")
	}

	account, err := db_organization.Query.GetAccount(ctx, pgtype.UUID{
		Bytes: u,
		Valid: true,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.GetAccountResponse{
		Account: &endpoints.Account{
			AccountId:         req.Msg.AccountId,
			UserId:            account.UserID,
			FirstName:         account.FirstName,
			LastName:          account.LastName,
			Phone:             account.Phone,
			Email:             account.Email.String,
			IsAdmin:           account.IsAdmin,
			IsCustomerSupport: account.IsCustomerSupport,
			IsDriver:          account.IsDriver,
		},
	}), nil
}
