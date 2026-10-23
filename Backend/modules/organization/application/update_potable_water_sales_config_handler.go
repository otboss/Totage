package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"

	"connectrpc.com/connect"
)

func (s *Handlers) UpdatePotableWaterSalesConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.UpdatePotableWaterSalesConfigCommand],
) (*connect.Response[endpoints.UpdatePotableWaterSalesConfigResponse], error) {

	return connect.NewResponse(&endpoints.UpdatePotableWaterSalesConfigResponse{}), nil
}
