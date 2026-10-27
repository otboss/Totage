package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	"waterordersystem/GOBACKEND/modules/organization/domain/templates"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	endpoints "waterordersystem/GOBACKEND/modules/payment/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Handlers) PaypalCaptureOrderHandler(
	ctx context.Context,
	req *connect.Request[endpoints.PaypalCaptureOrderCommand],
) (*connect.Response[endpoints.PaypalCaptureOrderResponse], error) {

	u, err := uuid.Parse(req.Msg.OrderId)
	if err != nil {
		return nil, errors.New("bad order id provided")
	}

	order, err := db_order.Query.GetOrderById(ctx, pgtype.UUID{
		Bytes: u,
		Valid: true,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	organizationPaypalCredentials, err := db_organization.Query.GetPayPalCredentialByOrganizationID(ctx, order.OrganizationID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	decryptedRestSdkClientID, err := utils.AES.Decrypt(utils.Env.ENC_KEY, organizationPaypalCredentials.RestSdkClientIDEncrypted.String)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}
	decryptedRestSdkSecret, err := utils.AES.Decrypt(utils.Env.ENC_KEY, organizationPaypalCredentials.RestSdkSecretEncrypted.String)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if order.PaypalOrderID.Valid == false {
		return nil, errors.New("paypal order not found")
	}

	organization, err := db_organization.Query.GetOrganization(ctx, order.OrganizationID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	paypalCredentials, err := db_organization.Query.GetPayPalCredentialByOrganizationID(ctx, organization.ID)
	if err == pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("organization is not subscribed"),
		)
	} else {
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		paypalSubscriptionDetails, err := clients.PaypalApi.ShowSubscriptionDetails(utils.Env.PAYPAL_REST_SDK_CLIENT_ID, utils.Env.PAYPAL_REST_SDK_CLIENT_SECRET, paypalCredentials.SubscriptionID)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		if paypalSubscriptionDetails.Status != clients.PaypalSubscriptionStatusActive {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("organization is not subscribed"),
			)
		}
	}

	captureOrderResponse, err := clients.PaypalApi.CaptureOrder(decryptedRestSdkClientID, decryptedRestSdkSecret, order.PaypalOrderID.String)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	_, err = db_order.Query.UpdateOrderStatusToPending(ctx, order.ID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	user, err := db_user.Query.GetUser(ctx, order.UserID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	organizationMailList, err := db_organization.Query.GetOrganizationMailListByOrganizationId(ctx, organization.ID)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
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
				"Order Placed",
				templates.OrderPlaced(user, order),
			)
			if err != nil {
				log.Println(err.Error())
			}

			err = clients.EmailApi.SendEmail(
				fmt.Sprintf("\"%s\" <%s>", utils.Env.SMTP_NAME, utils.Env.SMTP_USER),
				emailRecipients,
				"New Order Received",
				templates.NewOrderReceived(user, order),
			)
			if err != nil {
				log.Println(err.Error())
			}

			err = clients.TwilioApi.SendSms(
				user.Phone,
				fmt.Sprintf(
					"Hello customer. Your %s order with ID %s has been received and will be processed by %s. For more information please contact their customer support at %s",
					strings.Replace(string(order.Type), "_", " ", 0),
					order.ID.String(),
					organization.OrganizationName,
					organization.Phone,
				),
			)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	return connect.NewResponse(&endpoints.PaypalCaptureOrderResponse{
		Id:     captureOrderResponse.ID,
		Status: captureOrderResponse.Status,
		// TODO: complete remaining fields
	}), nil
}
