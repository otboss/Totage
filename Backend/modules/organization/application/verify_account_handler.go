package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) VerifyAccountHandler(
	ctx context.Context,
	req *connect.Request[endpoints.VerifyAccountCommand],
) (*connect.Response[endpoints.VerifyAccountResponse], error) {

	u, err := uuid.Parse(req.Msg.AccountId)
	if err != nil {
		return nil, errors.New("bad account id provided")
	}

	accountId := pgtype.UUID{
		Bytes: u,
		Valid: true,
	}

	_, err = db_organization.Query.GetAccount(ctx, accountId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	err = db_organization.Query.UpdateAccountInvitationAccepted(ctx, accountId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.VerifyAccountResponse{}), nil
}
