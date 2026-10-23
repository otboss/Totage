package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"

	"connectrpc.com/connect"
)

func (s *Handlers) CreateIceOrderHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateIceOrderCommand],
) (*connect.Response[endpoints.CreateOrderResponse], error) {
	// TODO:
	return connect.NewResponse(&endpoints.CreateOrderResponse{}), nil
}
