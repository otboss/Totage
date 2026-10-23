package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"

	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) GetPaypalOrderHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetPaypalOrderQuery],
) (*connect.Response[endpoints.GetPaypalOrderResponse], error) {

	// TODO:

	order, err := db_order.Query.GetOrderByPaypalOrderId(ctx, pgtype.Text{
		String: req.Msg.PaypalOrderId,
		Valid:  true,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.GetPaypalOrderResponse{
		// Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
		OrderId:     order.ID.String(),
		PaymentDate: order.PaymentDate.Int64,
		CreatedOn:   order.CreatedAt.Time.UnixMilli(),
	}), nil
}
