package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) ConfirmBankTransferHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ConfirmBankTransferCommand],
) (*connect.Response[endpoints.ConfirmBankTransferResponse], error) {

	decryptedOrderId, err := utils.AES.Decrypt(utils.Env.ENC_KEY, req.Msg.OrderIdEncrypted)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	u, err := uuid.Parse(decryptedOrderId)
	if err != nil {
		return nil, errors.New("bad session token provided")
	}

	_, err = db_order.Query.UpdateOrderStatusToPending(ctx, pgtype.UUID{
		Bytes: u,
		Valid: true,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.ConfirmBankTransferResponse{}), nil
}
