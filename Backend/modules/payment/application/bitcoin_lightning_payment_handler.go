package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/payment/api"

	"connectrpc.com/connect"
)

func (s *Handlers) BitcoinLightningPaymentHandler(
	ctx context.Context,
	req *connect.Request[endpoints.BitcoinLightningPaymentCommand],
) (*connect.Response[endpoints.BitcoinLightningPaymentResponse], error) {

	// TODO:implement

	return connect.NewResponse(&endpoints.BitcoinLightningPaymentResponse{
		// TODO: complete remaining fields
	}), nil
}
