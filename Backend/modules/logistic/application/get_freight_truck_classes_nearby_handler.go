package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	common "waterordersystem/GOBACKEND/shared/api"

	// shared "waterordersystem/GOBACKEND/shared/api"

	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"

	"connectrpc.com/connect"
)

func (s *Handlers) GetFreightTruckClassesNearbyHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetFreightTruckClassesNearbyQuery],
) (*connect.Response[endpoints.GetFreightTruckClassesNearbyResponse], error) {

	truckClasses, err := db_logistic.Query.GetTruckClassesNearby(ctx, db_logistic.GetTruckClassesNearbyParams{
		Longitude: req.Msg.PickupAddress.Longitude,
		Latitude:  req.Msg.PickupAddress.Latitude,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var response []*endpoints.FreightTruckClass

	for _, truckClass := range truckClasses {

		truckPropertyFromDb, err := db_logistic.Query.GetTruckPropertyByPropertyId(ctx, truckClass.Property)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		truckPropertyNum, exists := common.FreightTruckProperty_value[truckPropertyFromDb.Property]
		if exists == false {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		freightTruckClassNum, exists := common.FreightTruckClass_value[truckClass.Class.String]
		if exists == false {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		entry := endpoints.FreightTruckClass{
			Class:    common.FreightTruckClass(freightTruckClassNum),
			Property: common.FreightTruckProperty(truckPropertyNum),
		}
		response = append(response, &entry)
	}

	return connect.NewResponse(&endpoints.GetFreightTruckClassesNearbyResponse{
		FreightTruckClasses: response,
	}), nil
}
