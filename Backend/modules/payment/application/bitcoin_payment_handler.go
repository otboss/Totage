package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/payment/api"

	"connectrpc.com/connect"
)

func (s *Handlers) BitcoinPaymentHandler(
	ctx context.Context,
	req *connect.Request[endpoints.BitcoinPaymentCommand],
) (*connect.Response[endpoints.BitcoinPaymentResponse], error) {

	// TODO:implement

	return connect.NewResponse(&endpoints.BitcoinPaymentResponse{
		// TODO: complete remaining fields
	}), nil
}
