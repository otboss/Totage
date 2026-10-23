package application

import (
	"context"
	"sync"
	endpoints "waterordersystem/GOBACKEND/modules/logistic/api"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"

	"connectrpc.com/connect"
)

func (s *Handlers) ValidateAddressesHandler(
	ctx context.Context,
	req *connect.Request[endpoints.ValidateAddressesQuery],
) (*connect.Response[endpoints.ValidateAddressesResponse], error) {

	waypointAddresses := req.Msg.WaypointAddresses

	var wg sync.WaitGroup
	wg.Add(len(waypointAddresses))

	results := validationResult{}

	for i := 1; i < len(waypointAddresses); i++ {

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			var originCoordinates *clients.Coordinates
			var destinationCoordinates *clients.Coordinates

			if waypointAddresses[i-1].Coordinates != nil {
				originCoordinates = &clients.Coordinates{
					Latitude:  float64(waypointAddresses[i-1].Coordinates.Latitude),
					Longitude: float64(waypointAddresses[i-1].Coordinates.Longitude),
				}
			}

			if waypointAddresses[i].Coordinates != nil {
				destinationCoordinates = &clients.Coordinates{
					Latitude:  float64(waypointAddresses[i].Coordinates.Latitude),
					Longitude: float64(waypointAddresses[i].Coordinates.Longitude),
				}
			}

			distanceCalculation, err := clients.GoogleCloudApi.DistanceBetweenLocations(clients.Location{
				Address:     waypointAddresses[i-1].Address,
				Coordinates: originCoordinates,
			}, clients.Location{
				Address:     waypointAddresses[i].Address,
				Coordinates: destinationCoordinates,
			})
			if err != nil {
				results.AddValidation(int64(i), false)
				return
			}

			if distanceCalculation.Rows[0].Elements[0].Status == "ZERO_RESULTS" || distanceCalculation.Rows[0].Elements[0].Status == "NOT_FOUND" {
				results.AddValidation(int64(i), false)
				return
			}

			results.AddValidation(int64(i), true)
		}(
			&wg,
		)

	}

	wg.Wait()

	var validationResponse []*endpoints.ValidateAddressesResponse_ValidationResponse

	for _, result := range results.result {
		validationResponse = append(validationResponse, &endpoints.ValidateAddressesResponse_ValidationResponse{
			Index:              result.index,
			ValidationResponse: result.result,
		})
	}

	return connect.NewResponse(&endpoints.ValidateAddressesResponse{
		ValidationResponse: validationResponse,
	}), nil
}

type validationResult struct {
	mu     sync.Mutex
	result []struct {
		index  int64
		result bool
	}
}

func (t *validationResult) AddValidation(index int64, isValid bool) {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.result = append(t.result, struct {
		index  int64
		result bool
	}{
		index:  index,
		result: isValid,
	})
}
