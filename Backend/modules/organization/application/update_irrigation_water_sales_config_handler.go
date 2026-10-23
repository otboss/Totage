package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"

	"connectrpc.com/connect"
)

func (s *Handlers) UpdateIrrigationWaterSalesConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.UpdateIrrigationWaterSalesConfigCommand],
) (*connect.Response[endpoints.UpdateIrrigationWaterSalesConfigResponse], error) {

	return connect.NewResponse(&endpoints.UpdateIrrigationWaterSalesConfigResponse{}), nil
}
