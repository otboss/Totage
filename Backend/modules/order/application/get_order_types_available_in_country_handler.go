package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"

	"connectrpc.com/connect"
)

func (s *Handlers) GetOrderTypesAvailableInCountryHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetOrderTypesAvailableInCountryQuery],
) (*connect.Response[endpoints.GetOrderTypesAvailableInCountryResponse], error) {

	var originCountryCode string

	if req.Msg.Country == nil {
		originCountryCode = req.Header().Get("CF-IPCountry")
		if originCountryCode == "" {
			return nil, errors.New("origin country code header CF-IPCountry is missing")
		}
	} else {
		originCountryCode = req.Msg.Country.String()
	}

	supportedOrderTypes, err := db_organization.Query.GetOrganizationSupportedOrderTypes(ctx, db_organization.SharedCountryCode(originCountryCode))
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var supportedOrderTypesResponse []common.OrderType

	for _, orderType := range supportedOrderTypes {
		supportedOrderTypesResponse = append(supportedOrderTypesResponse, common.OrderType(common.OrderType_value[string(orderType)]))
	}

	return connect.NewResponse(&endpoints.GetOrderTypesAvailableInCountryResponse{
		SupportedOrderTypes: supportedOrderTypesResponse,
	}), nil
}
