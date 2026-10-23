package application

import (
	"context"
	"errors"
	"sync"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	// db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) GetFreightOrderTotalHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetFreightOrderTotalQuery],
) (*connect.Response[endpoints.GetOrderTotalResponse], error) {

	truck, err := db_logistic.Query.GetTruck(ctx, req.Msg.TruckId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	organization, err := db_organization.Query.GetOrganizationByTruckId(ctx, req.Msg.TruckId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if organization.IsApproved == false {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	org, err := db_organization.Query.GetOrganization(ctx, organization.ID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	err = HelperFunctions.CheckOrganizationSubscription(ctx, org)
	if err != nil {
		return nil, err
	}

	freightConfig, err := db_organization.Query.GetFreightConfig(ctx, organization.ID)
	if err == pgx.ErrNoRows || (err == nil && freightConfig.IsEnabled.Bool == false) {
		return nil, errors.New("order type is not enabled for organization")
	}
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	waypointAddresses := req.Msg.WaypointAddresses

	var wg sync.WaitGroup
	wg.Add(len(waypointAddresses))

	totalWaypointDistanceInKm := totalDistanceInKm{}
	addressError := false

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
				addressError = true
				return
			}

			if distanceCalculation.Rows[0].Elements[0].Status == "ZERO_RESULTS" || distanceCalculation.Rows[0].Elements[0].Status == "NOT_FOUND" {
				addressError = true
				return
			}

			totalWaypointDistanceInKm.Add(int32(distanceCalculation.Rows[0].Elements[0].Distance.Value / 1000))
		}(
			&wg,
		)

	}

	wg.Wait()

	if addressError == true {
		return nil, errors.New("invalid address provided")
	}

	total := HelperFunctions.CalculateFreightOrderTotal(truck, freightConfig, organization.DutyFreeKm, float64(totalWaypointDistanceInKm.total))

	response, err := HelperFunctions.BuildOrderTotalResponse(
		ctx,
		db_organization.SharedOrderTypeFREIGHT,
		db_organization.GetOrganizationRow{
			ID:                organization.ID,
			OrganizationName:  organization.OrganizationName,
			Phone:             organization.Phone,
			Email:             organization.Email,
			Country:           organization.Country,
			Currency:          organization.Currency,
			Address:           organization.Address,
			DutyFreeKm:        organization.DutyFreeKm,
			MaxOrderQueueSize: organization.MaxOrderQueueSize,
			AddressLongitude:  organization.AddressLongitude,
			AddressLatitude:   organization.AddressLatitude,
			IsApproved:        organization.IsApproved,
		},
		total,
		req.Msg.DiscountCode,
	)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(response), nil
}

//
//
//

func (h *helperFunctions) BuildOrderTotalResponse(ctx context.Context, orderType db_organization.SharedOrderType, organization db_organization.GetOrganizationRow, total float64, discountCode *string) (*endpoints.GetOrderTotalResponse, error) {
	var err error
	var discountValue float64 = 0

	if discountCode != nil {
		discountValue, err = HelperFunctions.GetDiscountValue(ctx, total, *discountCode)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
	}

	bankAccountInfo, err := db_organization.Query.GetAllBankAccountInfo(ctx, organization.ID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var resBankAccountInfo []*endpoints.BankAccountInfo
	for _, info := range bankAccountInfo {

		typeVal, exists := endpoints.BankAccountInfo_BankAccountType_value[string(info.Type)]
		if exists == false {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		ownershipVal, exists := endpoints.BankAccountInfo_BankAccountOwnership_value[string(info.Ownership)]
		if exists == false {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		resBankAccountInfo = append(resBankAccountInfo, &endpoints.BankAccountInfo{
			Id:         info.ID,
			Bank:       info.Bank,
			BankBranch: info.BankBranch,
			Type:       endpoints.BankAccountInfo_BankAccountType(typeVal),
			Ownership:  endpoints.BankAccountInfo_BankAccountOwnership(ownershipVal),
			Number:     info.AccountNumber,
		})

	}

	exchangeRateUsd, err := clients.CurrencyApi.GetExchangeRateUsd(string(organization.Currency))
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	// timeslots := endpoints.Timeslots{
	// 	MinSelectionTimeInMilliseconds: organization.MinSelectionTimeInMilliseconds,
	// 	Slots:                          []*endpoints.Timeslots_JobTimeslot{},
	// }

	// jobTimeslot, err := db_logistic.Query.GetTimeslot(ctx, db_organization.GetTimeslotParams{
	// 	OrganizationID: organization.ID,
	// 	OrderType:      orderType,
	// })
	// if err != nil {
	// 	return nil, connect.NewError(
	// 		connect.CodeInternal,
	// 		errors.New("internal server error"),
	// 	)
	// }

	// var weekdays []common.Weekdays
	// for _, weekday := range jobTimeslot.Weekdays {
	// 	weekdays = append(weekdays, common.Weekdays(weekday))
	// }

	paypalCredentials, err := db_organization.Query.GetPaypalCredentials(ctx, organization.ID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	var decryptedRestSdkClientID *string

	if paypalCredentials.RestSdkClientIDEncrypted.Valid == true {
		decryptionResult, err := utils.AES.Decrypt(utils.Env.ENC_KEY, paypalCredentials.RestSdkClientIDEncrypted.String)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
		decryptedRestSdkClientID = &decryptionResult
	}

	return &endpoints.GetOrderTotalResponse{
		Total: &endpoints.GetOrderTotalResponse_Total{
			Value:           total,
			DiscountValue:   discountValue,
			Currency:        string(organization.Currency),
			ExchangeRateUsd: exchangeRateUsd,
		},
		Organization: &endpoints.GetOrderTotalResponse_Organization{
			Id:              organization.ID,
			Name:            organization.OrganizationName,
			Address:         organization.Address,
			Phone:           organization.Phone,
			PaypalClientId:  decryptedRestSdkClientID,
			BankAccountInfo: resBankAccountInfo,
		},
	}, nil
}

func (h *helperFunctions) CalculateFreightOrderTotal(truck db_logistic.GetTruckRow, freightConfig db_organization.GetFreightConfigRow, dutyFreeKm float64, roundTripDistance float64) float64 {
	total := float64(freightConfig.FreightPriceMin)
	if float64(roundTripDistance) > dutyFreeKm {
		total += (float64(roundTripDistance) * float64(truck.FreightPricePerKm))
	}
	return total
}
