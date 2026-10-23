package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"

	"connectrpc.com/connect"
)

func (s *Handlers) GeocodingHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GeocodingQuery],
) (*connect.Response[endpoints.GeocodingResponse], error) {

	resp, err := clients.GoogleCloudApi.GeocodeAddress(req.Msg.Address)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.GeocodingResponse{
		Coordinates: &common.Coordinates{
			Latitude:  resp.Results[0].Geometry.Location.Lat,
			Longitude: resp.Results[0].Geometry.Location.Lng,
		},
	}), nil
}
