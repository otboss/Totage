package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"

	"connectrpc.com/connect"
)

func (s *Handlers) UpdateCharcoalSalesConfigHandler(
	ctx context.Context,
	req *connect.Request[endpoints.UpdateCharcoalSalesConfigCommand],
) (*connect.Response[endpoints.UpdateCharcoalSalesConfigResponse], error) {

	return connect.NewResponse(&endpoints.UpdateCharcoalSalesConfigResponse{}), nil
}
