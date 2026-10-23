package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	common "waterordersystem/GOBACKEND/shared/api"

	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"

	"connectrpc.com/connect"
)

func (s *Handlers) GetFreightTrucksNearbyHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetFreightTrucksNearbyQuery],
) (*connect.Response[endpoints.GetTrucksNearbyResponse], error) {

	trucks, err := db_logistic.Query.GetTrucksClosestToPoint(ctx, db_logistic.GetTrucksClosestToPointParams{
		Longitude:  req.Msg.PickupAddress.Longitude,
		Latitude:   req.Msg.PickupAddress.Latitude,
		OrderType:  db_logistic.SharedOrderTypeFREIGHT,
		Limitation: 10,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var trucksResponse []*endpoints.GetTrucksNearbyResponse_TruckResult
	var bedHeight float32

	for _, truck := range trucks {
		if truck.BedHeightMm.Valid != false {
			bedHeight = float32(truck.BedHeightMm.Float64)
		}

		truckBedSize := common.BedSize{
			LengthMm: float32(truck.BedLengthMm),
			WidthMm:  float32(truck.BedWidthMm),
			HeightMm: &bedHeight,
		}

		truckPropertyFromDb, err := db_logistic.Query.GetTruckPropertyByPropertyId(ctx, truck.Property)
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

		truckResponse := endpoints.GetTrucksNearbyResponse_TruckResult{
			DistanceAwayKm: truck.DistanceMeters / 1000,
			Truck: &common.Truck{
				Id:                truck.ID,
				VehicleId:         truck.VehicleID,
				BedSize:           &truckBedSize,
				Property:          common.FreightTruckProperty(truckPropertyNum),
				FreightPricePerKm: float32(truck.FreightPricePerKm),
				Cost:              float32(truck.Cost),
				// TODO:
				// Class: ,
			},
		}
		trucksResponse = append(trucksResponse, &truckResponse)
	}

	return connect.NewResponse(&endpoints.GetTrucksNearbyResponse{
		Trucks: trucksResponse,
	}), nil
}
