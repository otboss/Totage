package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handlers) GetOrganizationOrdersHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetOrganizationOrdersQuery],
) (*connect.Response[endpoints.GetOrganizationOrdersResponse], error) {

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
	if err == pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodePermissionDenied,
			errors.New("request not permitted"),
		)
	}
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	orders, err := db_order.Query.GetOrdersByOrganization(ctx, req.Msg.OrganizationId)
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

	return connect.NewResponse(&endpoints.GetOrganizationOrdersResponse{
		Orders: orderResponse,
	}), nil
}
