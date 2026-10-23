package application

import (
	"context"
	"errors"

	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) CompleteLogisticHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CompleteLogisticCommand],
) (*connect.Response[endpoints.CompleteLogisticResponse], error) {

	u, err := uuid.Parse(req.Msg.OrderId)
	if err != nil {
		return nil, errors.New("bad session token provided")
	}

	orderId := pgtype.UUID{
		Bytes: u,
		Valid: true,
	}

	var res *connect.Response[endpoints.CompleteLogisticResponse]
	_, err = db_order.Query.GetOrderById(ctx, orderId)

	if err != nil {
		res = connect.NewResponse(&endpoints.CompleteLogisticResponse{
			Status: endpoints.LogisticCompleteStatus_NOT_FOUND,
		})
		return res, nil
	}

	err = db_order.Query.UpdateOrderStatusToComplete(ctx, orderId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.CompleteLogisticResponse{
		Status: endpoints.LogisticCompleteStatus_COMPLETED,
	}), nil
}
