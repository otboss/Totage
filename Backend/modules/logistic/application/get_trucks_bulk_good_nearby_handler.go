package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"

	"connectrpc.com/connect"
)

func (s *Handlers) GetBulkGoodTrucksNearbyHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetBulkGoodTrucksNearbyQuery],
) (*connect.Response[endpoints.GetTrucksNearbyResponse], error) {

	var orderType db_logistic.SharedOrderType
	switch req.Msg.OrderType {
	case common.OrderType_BLOCK:
		orderType = db_logistic.SharedOrderTypeBLOCK
		break
	case common.OrderType_WATER_PURIFIED:
		orderType = db_logistic.SharedOrderTypeWATERPURIFIED
		break
	case common.OrderType_WATER_UNPURIFIED:
		orderType = db_logistic.SharedOrderTypeWATERUNPURIFIED
		break
	case common.OrderType_CHARCOAL:
		orderType = db_logistic.SharedOrderTypeCHARCOAL
		break
	case common.OrderType_EGG:
		orderType = db_logistic.SharedOrderTypeEGG
		break
	case common.OrderType_COCONUT:
		orderType = db_logistic.SharedOrderTypeCOCONUT
		break
	default:
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid order type specified"),
		)
	}

	trucks, err := db_logistic.Query.GetTrucksClosestToPoint(ctx, db_logistic.GetTrucksClosestToPointParams{
		Longitude:  req.Msg.Address.Longitude,
		Latitude:   req.Msg.Address.Latitude,
		OrderType:  orderType,
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
			DistanceAwayKm: truck.DistanceMeters,
			Truck: &common.Truck{
				Id:                truck.ID,
				VehicleId:         truck.VehicleID,
				BedSize:           &truckBedSize,
				Property:          common.FreightTruckProperty(truckPropertyNum),
				FreightPricePerKm: float32(truck.FreightPricePerKm),
				Cost:              float32(truck.Cost),
				IsPublic:          truck.IsPublic,
			},
		}
		trucksResponse = append(trucksResponse, &truckResponse)
	}

	return connect.NewResponse(&endpoints.GetTrucksNearbyResponse{
		Trucks: trucksResponse,
	}), nil
}
