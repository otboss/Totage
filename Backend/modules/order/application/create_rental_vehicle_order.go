package application

import (
	"context"
	"errors"
	"fmt"
	"time"
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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nyaruka/phonenumbers"
)

func (s *Handlers) CreateRentalVehicleOrderHandler(
	ctx context.Context,
	req *connect.Request[endpoints.CreateRentalVehicleOrderCommand],
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

	isRentalVehiclesEnabled, err := db_logistic.Query.GetRentalVehiclesEnabled(ctx, req.Msg.OrderDetails.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if isRentalVehiclesEnabled == false {
		return nil, errors.New("order type is not enabled for organization")
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

	scheduledRentals, err := db_order.Query.GetRentalVehicleByRentTimespan(ctx, db_order.GetRentalVehicleByRentTimespanParams{
		StartDate: req.Msg.OrderDetails.StartDate,
		EndDate:   req.Msg.OrderDetails.EndDate,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if len(scheduledRentals) > 0 {
		return nil, errors.New("vehicle is not available at the selected time")
	}

	rentalVehicle, err := db_logistic.Query.GetRentalVehicle(ctx, req.Msg.OrderDetails.RentalVehicleId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	user, err := db_user.Query.GetUser(ctx, claims.UserID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	total, err := HelperFunctions.CalculateRentalVehicleOrderTotal(
		db_logistic.LogisticRentalVehicle{
			MinDaysRental:   rentalVehicle.MinDaysRental,
			DailyRentalCost: rentalVehicle.DailyRentalCost,
		},
		time.UnixMilli(req.Msg.OrderDetails.StartDate), time.UnixMilli(req.Msg.OrderDetails.EndDate),
	)
	if err != nil {
		return nil, err
	}

	orderId, paypalResponseObj, err := HelperFunctions.FinalizePaymentAndSaveOrder(
		ctx,
		total,
		db_order.SharedOrderTypeVEHICLERENTAL,
		organization,
		user,
		req.Msg.OrderDetails.DiscountCode,
		req.Msg.OrderDetails.PaymentMethod,
		req.Msg.OrderDetails.BankAccountInfoId,
		[]*common.Coordinates{},
		0,
		&req.Msg.OrderDetails.StartDate,
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

	_, err = db_order.Query.CreateRentalVehicleOrder(ctx, db_order.CreateRentalVehicleOrderParams{
		OrderID: pgtype.UUID{
			Bytes: u,
			Valid: true,
		},
		RentalVehicleID: req.Msg.OrderDetails.RentalVehicleId,
		StartDate: pgtype.Int8{
			Int64: req.Msg.OrderDetails.StartDate,
			Valid: true,
		},
		EndDate: pgtype.Int8{
			Int64: req.Msg.OrderDetails.EndDate,
			Valid: true,
		},
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
