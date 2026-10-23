package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"

	"connectrpc.com/connect"
)

func (s *Handlers) CreateCharcoalOrderHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateCharcoalOrderCommand],
) (*connect.Response[endpoints.CreateOrderResponse], error) {
	// TODO:
	return connect.NewResponse(&endpoints.CreateOrderResponse{}), nil
}
