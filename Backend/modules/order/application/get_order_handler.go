package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handlers) GetOrderHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetOrderQuery],
) (*connect.Response[endpoints.GetOrderResponse], error) {

	u, err := uuid.Parse(req.Msg.OrderId)
	if err != nil {
		return nil, errors.New("bad order id provided")
	}

	orderId := pgtype.UUID{
		Bytes: u,
		Valid: true,
	}

	order, err := db_order.Query.GetOrder(ctx, orderId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	orderStatusVal, exists := endpoints.GetOrderResponse_OrderStatus_value[string(order.Status)]
	if exists == false {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	orderTypeVal, exists := common.OrderType_value[string(order.Type)]
	if exists == false {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var paymentDate *timestamppb.Timestamp

	if order.PaymentDate.Valid == true {
		paymentDate = timestamppb.New(utils.TimeFromMillis(order.PaymentDate.Int64))
	}

	var jobDate *int64
	if order.JobDate.Valid == true {
		jobDate = &order.JobDate.Int64
	}

	return connect.NewResponse(&endpoints.GetOrderResponse{
		PaymentDate:           paymentDate,
		OrderStatus:           endpoints.GetOrderResponse_OrderStatus(orderStatusVal),
		CreatedAt:             timestamppb.New(order.CreatedAt.Time),
		OrderType:             common.OrderType(orderTypeVal),
		JobDate:               jobDate,
		DistanceRoundTripInKm: order.DistanceRoundTripInKm,
		TotalInUsd:            order.TotalInUsd,
	}), nil
}
