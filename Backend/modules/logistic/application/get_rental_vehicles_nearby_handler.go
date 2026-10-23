package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"

	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"

	"connectrpc.com/connect"
)

func (s *Handlers) GetRentalVehiclesNearbyHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetRentalVehiclesNearbyQuery],
) (*connect.Response[endpoints.GetRentalVehiclesNearbyResponse], error) {

	geocodingResult, err := clients.GoogleCloudApi.GeocodeAddress(req.Msg.Address)
	if err != nil || len(geocodingResult.Results) == 0 {
		return nil, errors.New("invalid address provided")
	}

	rentalVehicles, err := db_logistic.Query.GetRentalVehiclesPublicClosestToPoint(ctx, db_logistic.GetRentalVehiclesPublicClosestToPointParams{
		Longitude:       geocodingResult.Results[0].Geometry.Location.Lng,
		Latitude:        geocodingResult.Results[0].Geometry.Location.Lat,
		SeatingCapacity: req.Msg.SeatingCapacity,
		Limitation:      10,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var rentalVehiclesResponse []*endpoints.GetRentalVehiclesNearbyResponse_RentalVehicleResult

	for _, rentalVehicle := range rentalVehicles {

		vehicle, err := db_logistic.Query.GetVehicle(ctx, rentalVehicle.VehicleID)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		rentalVehiclesResponse = append(rentalVehiclesResponse, &endpoints.GetRentalVehiclesNearbyResponse_RentalVehicleResult{
			Id:                   rentalVehicle.ID,
			VehicleId:            rentalVehicle.VehicleID,
			OrganizationId:       vehicle.OrganizationID,
			Year:                 int32(rentalVehicle.Year),
			Make:                 rentalVehicle.Make,
			Model:                rentalVehicle.Model,
			SeatingCapacity:      int32(rentalVehicle.MaxCapacity),
			FuelType:             string(rentalVehicle.FuelType),
			DailyRentalCost:      rentalVehicle.DailyRentalCost,
			MinDaysRental:        int32(rentalVehicle.MinDaysRental),
			IsManualTransmission: rentalVehicle.IsManualTransmission,
		})
	}

	return connect.NewResponse(&endpoints.GetRentalVehiclesNearbyResponse{
		RentalVehicles: rentalVehiclesResponse,
	}), nil
}
