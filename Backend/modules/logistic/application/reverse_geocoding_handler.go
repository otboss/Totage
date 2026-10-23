package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"

	"connectrpc.com/connect"
)

func (s *Handlers) ReverseGeocodingHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ReverseGeocodingQuery],
) (*connect.Response[endpoints.ReverseGeocodingResponse], error) {

	address, err := clients.GoogleCloudApi.ReverseGeocoding(req.Msg.Coordinates.Latitude, req.Msg.Coordinates.Longitude)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if len(address.Results) == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("unable to process coordinates"),
		)
	}

	return connect.NewResponse(&endpoints.ReverseGeocodingResponse{
		Address: address.Results[0].FormattedAddress,
	}), nil
}
