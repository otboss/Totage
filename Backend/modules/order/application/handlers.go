package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	db_logistic "waterordersystem/GOBACKEND/modules/logistic/infrastructure"
	endpoints "waterordersystem/GOBACKEND/modules/order/api"
	"waterordersystem/GOBACKEND/modules/order/api/endpointsconnect"
	"waterordersystem/GOBACKEND/modules/organization/domain/templates"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/shared/infrastructure/middleware"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	api "waterordersystem/GOBACKEND/modules/order/api"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
)

func init() {
	handlers := &Handlers{}

	var path string
	var handler http.Handler

	path, handler = endpointsconnect.NewOrderServiceHandler(
		handlers,
		connect.WithInterceptors(
			middleware.AuthenticationInterceptor(),
			middleware.ValidationInterceptor(),
		),
	)
	utils.Mux.Handle(path, utils.RecoverMiddleware(handler))
	fmt.Println(path)
}

type Handlers struct {
}

//
//
//
//
//
//
//

type totalDistanceInKm struct {
	mu    sync.Mutex
	total int32
}

func (t *totalDistanceInKm) Add(distance int32) {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.total += distance
}

type helperFunctions struct{}

var HelperFunctions = helperFunctions{}

// Checks if any trucks for the organization is available for the job
func (h *helperFunctions) CheckTruckAvailability(ctx context.Context, organization db_organization.GetOrganizationRow, orderType db_organization.SharedOrderType, jobDate int64) error {

	if jobDate < time.Now().UnixMilli() {
		return errors.New("invalid job date selected")
	}

	// timeslot, err := db_organization.Query.GetTimeslot(ctx, db_organization.GetTimeslotParams{
	// 	OrganizationID: organization.ID,
	// 	OrderType:      orderType,
	// })
	// if err != nil {
	// 	return connect.NewError(
	// 		connect.CodeInternal,
	// 		errors.New("internal server error"),
	// 	)
	// }

	// if jobDate-time.Now().UnixMilli() <= timeslot.AverageJobDuration {
	// 	return errors.New("invalid job date selected")
	// }

	if organization.MaxOrderQueueSize.Valid == true {
		upcomingOrderCount, err := db_order.Query.GetIncompleteOrderCountByOrganization(
			ctx,
			organization.ID,
		)
		if err != nil {
			return errors.New("internal server error")
		}
		if upcomingOrderCount >= int32(organization.MaxOrderQueueSize.Int16) {
			return errors.New("organization order queue is full")
		}
	}

	upcomingOrders, err := db_order.Query.GetOrdersIncompletedForMonth(ctx, db_order.GetOrdersIncompletedForMonthParams{
		OrganizationID: organization.ID,
		JobMonth: pgtype.Int2{
			Int16: int16(utils.TimeFromMillis(jobDate).Month()) - 1,
			Valid: true,
		},
		JobYear: pgtype.Int2{
			Int16: int16(utils.TimeFromMillis(jobDate).Year()),
			Valid: true,
		},
	})
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	schedule := make(map[int64][]db_order.OrderOrder)
	for _, order := range upcomingOrders {
		if order.JobDate.Valid == false {
			continue
		}

		// TODO: Look into this
		// schedule[order.JobDate.Int64] = append(schedule[order.JobDate.Int64], order)
	}

	truckCount, err := db_logistic.Query.GetTruckCount(ctx, organization.ID)
	if err != nil {
		return connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, exists := schedule[jobDate]
	if exists {
		if int64(len(schedule[jobDate])) >= truckCount {
			return errors.New("no trucks available for job")
		}
	}

	return nil
}

func (h *helperFunctions) CheckOrganizationSubscription(ctx context.Context, organization db_organization.GetOrganizationRow) error {

	if utils.Env.IS_ORGANIZATION_SUBSCRIPTION_ENABLED == false {
		return nil
	}

	paypalCredentials, err := db_organization.Query.GetPaypalCredentials(ctx, organization.ID)
	if err == pgx.ErrNoRows {
		return errors.New("organization is not subscribed")
	} else {
		if err != nil {
			return errors.New("internal server error")
		}

		paypalSubscriptionDetails, err := clients.PaypalApi.ShowSubscriptionDetails(utils.Env.PAYPAL_REST_SDK_CLIENT_ID, utils.Env.PAYPAL_REST_SDK_CLIENT_SECRET, paypalCredentials.SubscriptionID)
		if err != nil {
			return errors.New("internal server error")
		}

		if paypalSubscriptionDetails.Status != clients.PaypalSubscriptionStatusActive {
			return errors.New("organization is not subscribed")
		}
	}

	return nil
}

type OrderId string

func (h *helperFunctions) FinalizePaymentAndSaveOrder(
	ctx context.Context,
	total float64,
	orderType db_order.SharedOrderType,
	organization db_organization.GetOrganizationRow,
	user db_user.GetUserRow,
	discountCode *string,
	paymentMethod endpoints.PaymentMethod,
	bankAccountInfoId *int64,
	waypoints []*common.Coordinates,
	roundTripDistance float64,
	jobDate *int64,
) (OrderId, *endpoints.CreateOrderResponse_PaypalOrderInfo, error) {
	var orderId OrderId
	var paypalResponseObj *endpoints.CreateOrderResponse_PaypalOrderInfo
	if string(paymentMethod) == api.PaymentMethod_name[int32(api.PaymentMethod_BANK_TRANSFER)] {
		if user.Email.Valid == false {
			return orderId, paypalResponseObj, errors.New("email address required for bank transfer")
		}
	}

	exchangeRate, err := clients.CurrencyApi.GetExchangeRateUsd(string(organization.Currency))
	if err != nil {
		return orderId, paypalResponseObj, errors.New("internal server error")
	}

	if discountCode != nil {
		total, err = h.GetDiscountValue(ctx, total, *discountCode)
		if err != nil {
			return orderId, paypalResponseObj, errors.New("internal server error")
		}
	}

	var paypalCreateOrderResponse *clients.PayPalCreateOrderResponse
	orderStatus := db_order.OrderOrderStatusPENDING

	switch paymentMethod {
	case endpoints.PaymentMethod_PAYPAL:

		paypalCredentials, err := db_organization.Query.GetPaypalCredentials(ctx, organization.ID)
		if err != nil {
			return orderId, paypalResponseObj, errors.New("internal server error")
		}

		if paypalCredentials.RestSdkClientIDEncrypted.Valid == false || paypalCredentials.RestSdkSecretEncrypted.Valid == false {
			return orderId, paypalResponseObj, errors.New("paypal not supported")
		}

		decryptedRestSdkClientID, err := utils.AES.Decrypt(utils.Env.ENC_KEY, paypalCredentials.RestSdkClientIDEncrypted.String)
		if err != nil {
			return orderId, paypalResponseObj, errors.New("internal server error")
		}
		decryptedRestSdkSecret, err := utils.AES.Decrypt(utils.Env.ENC_KEY, paypalCredentials.RestSdkSecretEncrypted.String)
		if err != nil {
			return orderId, paypalResponseObj, errors.New("internal server error")
		}

		paypalResponse, err := clients.PaypalApi.CreateOrder(decryptedRestSdkClientID, decryptedRestSdkSecret, total/exchangeRate)
		if err != nil {
			return orderId, paypalResponseObj, errors.New("internal server error")
		}

		paypalCreateOrderResponse = &paypalResponse

		var links []*api.CreateOrderResponse_PaypalOrderInfo_Link
		for _, link := range paypalResponse.Links {
			linkObj := api.CreateOrderResponse_PaypalOrderInfo_Link{
				Href:   link.Href,
				Rel:    link.Rel,
				Method: link.Method,
			}
			links = append(links, &linkObj)
		}

		paypalResponseObj = &endpoints.CreateOrderResponse_PaypalOrderInfo{
			Status: paypalResponse.Status,
			Links:  links,
		}
	case endpoints.PaymentMethod_BANK_TRANSFER:
		orderStatus = db_order.OrderOrderStatusUNPAID

		organizationMailList, err := db_organization.Query.GetOrganizationMailListByOrganizationId(ctx, organization.ID)
		if err != nil {
			return orderId, paypalResponseObj, errors.New("internal server error")
		}

		var emailRecipients []string
		for _, email := range organizationMailList {
			emailRecipients = append(emailRecipients, email.Email)
		}

		go func() {
			if utils.Env.EMAIL_DISABLED == false {
				err = clients.EmailApi.SendEmail(
					fmt.Sprintf("\"%s\" <%s>", utils.Env.SMTP_NAME, utils.Env.SMTP_USER),
					[]string{user.Email.String},
					"Bank Transfer Opened",
					templates.BankTransferOpened(),
				)
				if err != nil {
					log.Println(err.Error())
				}

				err = clients.EmailApi.SendEmail(
					fmt.Sprintf("\"%s\" <%s>", utils.Env.SMTP_NAME, utils.Env.SMTP_USER),
					emailRecipients,
					"New Bank Transfer Opened",
					templates.NewBankTransferOpened(),
				)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}()

	case endpoints.PaymentMethod_BITCOIN:
		// TODO: Handle bitcoin payment
	case endpoints.PaymentMethod_BITCOIN_LIGHTNING:
	}

	var paypalOrderId string
	if paypalCreateOrderResponse != nil {
		paypalOrderId = paypalCreateOrderResponse.ID
	}

	var bankAccountInfoIdNum int64
	if bankAccountInfoId != nil {
		bankAccountInfoIdNum = *bankAccountInfoId
	}

	var jDate int64
	if jobDate != nil {
		jDate = *jobDate
	}

	savedOrder, err := db_order.Query.CreateOrder(ctx, db_order.CreateOrderParams{
		PaymentMethod: db_order.OrderPaymentMethod(paymentMethod),
		PaypalOrderID: pgtype.Text{
			String: paypalOrderId,
			Valid:  paypalCreateOrderResponse != nil,
		},
		BankAccountID: pgtype.Int8{
			Int64: bankAccountInfoIdNum,
			Valid: bankAccountInfoId != nil,
		},
		UserID:         user.ID,
		OrganizationID: organization.ID,
		PaymentDate: pgtype.Int8{
			Int64: time.Now().UnixMilli(),
			Valid: true,
		},
		JobDate: pgtype.Int8{
			Int64: jDate,
			Valid: jobDate != nil,
		},
		JobMonth: pgtype.Int2{
			Int16: int16(utils.TimeFromMillis(jDate).Month()) - 1,
			Valid: jobDate != nil,
		},
		JobYear: pgtype.Int2{
			Int16: int16(utils.TimeFromMillis(jDate).Year()),
			Valid: jobDate != nil,
		},
		DistanceRoundTripInKm: roundTripDistance,
		Type:                  orderType,
		Status:                orderStatus,
		TotalInUsd:            total,
	})
	if err != nil {
		return orderId, paypalResponseObj, errors.New("internal server error")
	}

	for _, waypoint := range waypoints {
		u, err := uuid.Parse(string(orderId))
		if err != nil {
			return orderId, paypalResponseObj, errors.New("internal server error")
		}

		_, err = db_order.Query.CreateWaypoint(ctx, db_order.CreateWaypointParams{
			OrderID: pgtype.UUID{
				Bytes: u,
				Valid: true,
			},
			Longitude: waypoint.Longitude,
			Latitude:  waypoint.Latitude,
		})
		if err != nil {
			return orderId, paypalResponseObj, errors.New("internal server error")
		}
	}

	return OrderId(savedOrder.ID.String()), paypalResponseObj, nil
}

func (h *helperFunctions) GetDiscountValue(ctx context.Context, total float64, discountCode string) (float64, error) {
	discount, err := db_order.Query.GetDiscountCode(ctx, discountCode)
	if err != nil {
		return 0, err
	}

	isDiscountCodeValid, err := db_order.Query.IsDiscountCodeValid(ctx, discountCode)
	if err != nil {
		return 0, errors.New("invalid discount code")
	}

	if isDiscountCodeValid.Bool == true {
		switch discount.DiscountCodeType {
		case db_order.OrderDisconutCodeTypeFIXED:
			total -= discount.Value
		case db_order.OrderDisconutCodeTypePERCENT:
			total *= discount.Value / 100
		}
	}

	return total, nil
}
