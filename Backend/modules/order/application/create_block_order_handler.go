package application

import (
	"context"
	"errors"
	"fmt"
	api "waterordersystem/GOBACKEND/modules/order/api"
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

func (s *Handlers) CreateBlockOrderHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateBlockOrderCommand],
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

	organization, err := db_organization.Query.GetOrganization(ctx, int64(req.Msg.OrderDetails.OrganizationId))
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

	blockConfig, err := db_organization.Query.GetBlockSalesConfig(ctx, organization.ID)
	if err == pgx.ErrNoRows || (err == nil && blockConfig.IsEnabled.Bool == false) {
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

	// db_organization.Query.GetTimeslot(ctx, db_organization.GetTimeslotParams{})

	if req.Msg.OrderDetails.JobDate != nil {
		err = HelperFunctions.CheckTruckAvailability(ctx, organization, db_organization.SharedOrderTypeBLOCK, *req.Msg.OrderDetails.JobDate)
		if err != nil {
			return nil, err
		}
	}

	destinationLocation := clients.Location{
		Address:     req.Msg.OrderDetails.DeliveryLocation.Address,
		Coordinates: nil,
	}
	if req.Msg.OrderDetails.DeliveryLocation.Coordinates != nil {
		destinationLocation.Coordinates = &clients.Coordinates{
			Latitude:  float64(req.Msg.OrderDetails.DeliveryLocation.Coordinates.Latitude),
			Longitude: float64(req.Msg.OrderDetails.DeliveryLocation.Coordinates.Longitude),
		}
	}

	distanceCalculation, err := clients.GoogleCloudApi.DistanceBetweenLocations(
		clients.Location{
			Address:     organization.Address,
			Coordinates: nil,
		},
		destinationLocation,
	)

	if distanceCalculation.Rows[0].Elements[0].Status == "ZERO_RESULTS" || distanceCalculation.Rows[0].Elements[0].Status == "NOT_FOUND" {
		return nil, errors.New("an invalid address was provided")
	}

	user, err := db_user.Query.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if string(req.Msg.OrderDetails.PaymentMethod) == api.PaymentMethod_name[int32(api.PaymentMethod_BANK_TRANSFER)] {
		if req.Msg.Customer.Email == nil || user.Email.Valid == false {
			return nil, errors.New("email address required for bank transfer")
		}
	}

	jobRoundTripDistance := float64(distanceCalculation.Rows[0].Elements[0].Distance.Value/1000) * 2

	total := float64(req.Msg.OrderDetails.QuantityOfBlocks)*blockConfig.BlockPrice + blockConfig.BlockPriceMin
	if jobRoundTripDistance > organization.DutyFreeKm {
		total += jobRoundTripDistance * blockConfig.BlockPricePerKm
	}

	var waypointCoordinates []*common.Coordinates
	waypointCoordinates = append(waypointCoordinates, &common.Coordinates{
		Longitude: organization.AddressLongitude,
		Latitude:  organization.AddressLatitude,
	})

	if req.Msg.OrderDetails.DeliveryLocation.Coordinates != nil {
		waypointCoordinates = append(waypointCoordinates, &common.Coordinates{
			Longitude: req.Msg.OrderDetails.DeliveryLocation.Coordinates.Longitude,
			Latitude:  req.Msg.OrderDetails.DeliveryLocation.Coordinates.Latitude,
		})
	} else {
		geocodeResponse, err := clients.GoogleCloudApi.GeocodeAddress(req.Msg.OrderDetails.DeliveryLocation.Address)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
		if len(geocodeResponse.Results) == 0 {
			return nil, errors.New("bad delivery address provided")
		}
		waypointCoordinates = append(waypointCoordinates, &common.Coordinates{
			Longitude: geocodeResponse.Results[0].Geometry.Location.Lng,
			Latitude:  geocodeResponse.Results[0].Geometry.Location.Lat,
		})
	}

	orderId, paypalResponseObj, err := HelperFunctions.FinalizePaymentAndSaveOrder(
		ctx,
		total,
		db_order.SharedOrderTypeBLOCK,
		organization,
		user,
		req.Msg.OrderDetails.DiscountCode,
		req.Msg.OrderDetails.PaymentMethod,
		req.Msg.OrderDetails.BankAccountInfoId,
		waypointCoordinates,
		jobRoundTripDistance,
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

	_, err = db_order.Query.CreateBlockOrder(ctx, db_order.CreateBlockOrderParams{
		OrderID: pgtype.UUID{
			Bytes: u,
			Valid: true,
		},
		QuantityOfBlocks: int16(req.Msg.OrderDetails.QuantityOfBlocks),
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.CreateOrderResponse{
		OrderId:         string(orderId),
		PaypalOrderInfo: paypalResponseObj,
	}), nil
}
