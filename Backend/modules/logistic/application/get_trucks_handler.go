package application

import (
	"context"
	"errors"
	"sync"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"

	"connectrpc.com/connect"
)

func (s *Handlers) GetTrucksHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetTrucksQuery],
) (*connect.Response[endpoints.GetTrucksResponse], error) {

	var offsetId int64 = 0
	if req.Msg.Pagination.OffsetId != nil {
		offsetId = *req.Msg.Pagination.OffsetId
	}

	trucks, err := db_logistic.Query.GetTrucksByOrganization(ctx, db_logistic.GetTrucksByOrganizationParams{
		OrganizationID: req.Msg.OrganizationId,
		OffsetID:       offsetId,
		Limitation:     int32(req.Msg.Pagination.Limit),
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var trucksResponse []*endpoints.GetTrucksResponse_Truck
	var bedHeight float32

	var wg sync.WaitGroup
	wg.Add(len(trucks))
	var wgErrors []error

	for _, truck := range trucks {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			if truck.BedHeightMm.Valid != false {
				bedHeight = float32(truck.BedHeightMm.Float64)
			}

			truckBedSize := common.BedSize{
				LengthMm: float32(truck.BedLengthMm),
				WidthMm:  float32(truck.BedWidthMm),
				HeightMm: &bedHeight,
			}

			vehicle, err := db_logistic.Query.GetVehicle(ctx, truck.VehicleID)
			if err != nil {
				wgErrors = append(wgErrors, err)
				return
			}

			truckPropertyFromDb, err := db_logistic.Query.GetTruckPropertyByPropertyId(ctx, truck.Property)
			if err != nil {
				wgErrors = append(wgErrors, err)
				return
			}

			truckPropertyNum, exists := common.FreightTruckProperty_value[truckPropertyFromDb.Property]
			if exists == false {
				wgErrors = append(wgErrors, err)
				return
			}

			trucksResponse = append(trucksResponse, &endpoints.GetTrucksResponse_Truck{
				Truck: &common.Truck{
					Id:                truck.ID,
					VehicleId:         truck.VehicleID,
					BedSize:           &truckBedSize,
					Property:          common.FreightTruckProperty(truckPropertyNum),
					FreightPricePerKm: float32(truck.FreightPricePerKm),
					Cost:              float32(truck.Cost),
					IsPublic:          truck.IsPublic,
				},
				ChassisNumber: vehicle.ChassisNumber,
			})
		}(
			&wg,
		)
	}

	if len(wgErrors) > 0 {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.GetTrucksResponse{
		Trucks: trucksResponse,
	}), nil
}
