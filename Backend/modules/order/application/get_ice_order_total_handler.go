package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"

	"connectrpc.com/connect"
)

func (s *Handlers) GetIceOrderTotalHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetIceOrderTotalQuery],
) (*connect.Response[endpoints.GetOrderTotalResponse], error) {

	// TODO: Implement
	return connect.NewResponse(&endpoints.GetOrderTotalResponse{}), nil
}
