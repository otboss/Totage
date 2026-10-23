package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"

	"connectrpc.com/connect"
)

func (s *Handlers) GetRentalVehiclesHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetRentalVehiclesQuery],
) (*connect.Response[endpoints.GetRentalVehiclesResponse], error) {

	var offsetId int64 = 0
	if req.Msg.Pagination.OffsetId != nil {
		offsetId = *req.Msg.Pagination.OffsetId
	}

	rentalVehicles, err := db_logistic.Query.GetRentalVehiclesByOrganizationId(ctx, db_logistic.GetRentalVehiclesByOrganizationIdParams{
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

	var rentalVehiclesResponse []*endpoints.GetRentalVehiclesResponse_RentalVehicle

	for _, rentalVehicle := range rentalVehicles {
		rentalVehiclesResponse = append(rentalVehiclesResponse, &endpoints.GetRentalVehiclesResponse_RentalVehicle{
			RentalVehicleId: rentalVehicle.ID,
			Year:            int64(rentalVehicle.Year),
			Make:            rentalVehicle.Make,
			Model:           rentalVehicle.Model,
			// SeatingCapacity: int32(rentalVehicle.SeatingCapacity),
			// FuelType:        string(rentalVehicle.FuelType),
			// DailyRentalCost: float32(rentalVehicle.DailyRentalCost),
			// MinDaysRental:   int32(rentalVehicle.DailyRentalCost),
			// IsPublic:        rentalVehicle.IsPublic,
		})
	}

	return connect.NewResponse(&endpoints.GetRentalVehiclesResponse{
		RentalVehicles: rentalVehiclesResponse,
	}), nil
}
