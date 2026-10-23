package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handlers) GetUserOrdersHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetUserOrdersQuery],
) (*connect.Response[endpoints.GetUserOrdersResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	orders, err := db_order.Query.GetOrdersByUserId(ctx, claims.UserID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var orderResponse []*endpoints.GetOrderResponse

	for _, order := range orders {
		var paymentDate *timestamppb.Timestamp

		if order.PaymentDate.Valid == true {
			paymentDate = timestamppb.New(utils.TimeFromMillis(order.PaymentDate.Int64))
		}
		var createdAtDate *timestamppb.Timestamp

		if order.CreatedAt.Valid == true {
			createdAtDate = timestamppb.New(utils.TimeFromMillis(order.CreatedAt.Time.UnixMilli()))
		}

		var jobDate *int64
		if order.JobDate.Valid == true {
			jobDate = &order.JobDate.Int64
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

		orderResponse = append(orderResponse, &endpoints.GetOrderResponse{
			OrderId:               order.ID.String(),
			PaymentDate:           paymentDate,
			OrderStatus:           endpoints.GetOrderResponse_OrderStatus(orderStatusVal),
			CreatedAt:             createdAtDate,
			OrderType:             common.OrderType(orderTypeVal),
			JobDate:               jobDate,
			DistanceRoundTripInKm: order.DistanceRoundTripInKm,
			TotalInUsd:            order.TotalInUsd,
		})
	}

	return connect.NewResponse(&endpoints.GetUserOrdersResponse{
		Orders: orderResponse,
	}), nil
}
