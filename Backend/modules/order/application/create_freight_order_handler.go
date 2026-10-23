package application

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nyaruka/phonenumbers"
)

func (s *Handlers) CreateFreightOrderHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateFreightOrderCommand],
) (*connect.Response[endpoints.CreateOrderResponse], error) {

	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	countryCode, exist := common.CountryCodeAlpha2_name[int32(req.Msg.Customer.PhoneCountryCodeAlpha2)]
	if exist == false {
		return nil, fmt.Errorf("invalid country code provided")
	}

	num, err := phonenumbers.Parse(req.Msg.Customer.Phone, countryCode)
	if err != nil {
		return nil, errors.New("invalid phone number provided")
	}

	if !phonenumbers.IsValidNumber(num) {
		return nil, errors.New("invalid phone number provided")
	}

	if req.Msg.OrderDetails.PaymentMethod == endpoints.PaymentMethod_BANK_TRANSFER {
		if req.Msg.Customer.Email == nil {
			return nil, errors.New("email address is required for bank transfers")
		}
		if req.Msg.OrderDetails.BankAccountInfoId == nil {
			return nil, errors.New("bank account info is required for bank transfers")
		}
	}

	if len(req.Msg.OrderDetails.WaypointAddresses) < 2 {
		return nil, errors.New("at least two addresses should be provided")
	}

	if utils.Env.CAPTCHA_DISABLED == false {
		recaptchaValidation, err := clients.GoogleCloudApi.RecaptchaValidation(clients.RecaptchaEventRequest{
			Event: clients.RecaptchaEvent{
				Token:          req.Msg.RecaptchaEvent.Token,
				ExpectedAction: req.Msg.RecaptchaEvent.ExpectedAction,
				SiteKey:        req.Msg.RecaptchaEvent.SiteKey,
			},
		})
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		if recaptchaValidation.TokenProperties.Valid == false {
			return nil, errors.New("recaptcha validation failed")
		}
	}

	organization, err := db_organization.Query.GetOrganization(ctx, req.Msg.OrderDetails.OrganizationId)
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

	err = HelperFunctions.CheckOrganizationSubscription(ctx, organization)
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

	if organization.MaxOrderQueueSize.Valid == true {
		incompleteOrderCount, err := db_order.Query.GetIncompleteOrderCountByOrganization(ctx, req.Msg.OrderDetails.OrganizationId)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		if incompleteOrderCount >= int32(organization.MaxOrderQueueSize.Int16) {
			return nil, errors.New("organization order queue is full")
		}
	}

	if req.Msg.OrderDetails.JobDate != nil {
		err = HelperFunctions.CheckTruckAvailability(ctx, organization, db_organization.SharedOrderTypeFREIGHT, *req.Msg.OrderDetails.JobDate)
		if err != nil {
			return nil, err
		}
	}

	waypointAddresses := req.Msg.OrderDetails.WaypointAddresses
	waypointAddresses = append(waypointAddresses, waypointAddresses[0]) // round trip distance

	var waypoints []*waypointWithIndex
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
				waypoints = append(waypoints, &waypointWithIndex{
					Error: &err,
				})
				addressError = true
				return
			}

			if distanceCalculation.Rows[0].Elements[0].Status == "ZERO_RESULTS" || distanceCalculation.Rows[0].Elements[0].Status == "NOT_FOUND" {
				err := errors.New("invalid address provided")
				waypoints = append(waypoints, &waypointWithIndex{
					Error: &err,
				})
				addressError = true
				return
			}

			totalWaypointDistanceInKm.Add(int32(distanceCalculation.Rows[0].Elements[0].Distance.Value / 1000))

			x := i
			if i == 1 {
				x = 0
			}

			for true {
				var locationData clients.Coordinates

				if waypointAddresses[x].Coordinates != nil {
					locationData.Longitude = float64(waypointAddresses[x].Coordinates.Longitude)
					locationData.Latitude = float64(waypointAddresses[x].Coordinates.Latitude)
				} else {
					geocodeResult, err := clients.GoogleCloudApi.GeocodeAddress(waypointAddresses[x].Address)
					if err != nil {
						addressError = true
					} else {
						locationData.Longitude = geocodeResult.Results[0].Geometry.Location.Lng
						locationData.Latitude = geocodeResult.Results[0].Geometry.Location.Lat
					}
				}

				waypoints = append(waypoints, &waypointWithIndex{
					Index: int64(x),
					Waypoint: common.Coordinates{
						Longitude: locationData.Longitude,
						Latitude:  locationData.Latitude,
					},
					Address: waypointAddresses[x].Address,
				})

				if x == 0 {
					x = i
					continue
				}
				break
			}
		}(
			&wg,
		)

	}

	wg.Wait()

	if addressError == true {
		return nil, errors.New("invalid address provided")
	}

	user, err := db_user.Query.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	truck, err := db_logistic.Query.GetTruck(ctx, int64(req.Msg.OrderDetails.Truck.Id))
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	total := HelperFunctions.CalculateFreightOrderTotal(truck, freightConfig, organization.DutyFreeKm, float64(totalWaypointDistanceInKm.total))

	var waypointCoordinates []*common.Coordinates

	sort.Slice(waypoints, func(i, j int) bool {
		return waypoints[i].Index < waypoints[j].Index
	})

	for _, waypoint := range waypoints {
		waypointCoordinates = append(waypointCoordinates, &waypoint.Waypoint)
	}

	orderId, paypalResponseObj, err := HelperFunctions.FinalizePaymentAndSaveOrder(
		ctx,
		total,
		db_order.SharedOrderTypeFREIGHT,
		organization,
		user,
		req.Msg.OrderDetails.DiscountCode,
		req.Msg.OrderDetails.PaymentMethod,
		req.Msg.OrderDetails.BankAccountInfoId,
		waypointCoordinates[0:len(waypointCoordinates)-2], // Ignore last waypoint (ie. round trip coordinate)
		float64(totalWaypointDistanceInKm.total),
		req.Msg.OrderDetails.JobDate,
	)
	if err != nil {
		return nil, err
	}

	u, err := uuid.Parse(string(orderId))
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, err = db_order.Query.CreateFreightOrder(ctx, db_order.CreateFreightOrderParams{
		TruckID: req.Msg.OrderDetails.Truck.Id,
		OrderID: pgtype.UUID{
			Bytes: u,
			Valid: true,
		},
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	res := connect.NewResponse(&endpoints.CreateOrderResponse{
		OrderId:         string(orderId),
		PaypalOrderInfo: paypalResponseObj,
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

type waypointWithIndex struct {
	Index    int64
	Waypoint common.Coordinates
	Address  string
	Error    *error
}
