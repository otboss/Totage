package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"

	"connectrpc.com/connect"
)

func (s *Handlers) GetLogisticQueueSizeHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetLogisticQueueSizeQuery],
) (*connect.Response[endpoints.GetLogisticQueueSizeResponse], error) {

	count, err := db_order.Query.GetIncompleteOrderCountByOrganization(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}
	limit, err := db_organization.Query.GetOrganizationMaxOrderQueueLength(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var limitConverted *int32
	var limitRaw int32

	if limit.Valid {
		limitRaw = int32(limit.Int16)
		limitConverted = &limitRaw
	}

	return connect.NewResponse(&endpoints.GetLogisticQueueSizeResponse{
		CurrentSize: count,
		Limit:       limitConverted,
	}), nil
}
