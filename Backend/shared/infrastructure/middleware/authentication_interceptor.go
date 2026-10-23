package middleware

import (
	"context"
	"errors"
	"time"
	logistic_endpointsconnect "waterordersystem/GOBACKEND/modules/logistic/api/endpointsconnect"
	order_endpointsconnect "waterordersystem/GOBACKEND/modules/order/api/endpointsconnect"
	organization_endpointsconnect "waterordersystem/GOBACKEND/modules/organization/api/endpointsconnect"
	payment_endpointsconnect "waterordersystem/GOBACKEND/modules/payment/api/endpointsconnect"
	user_endpointsconnect "waterordersystem/GOBACKEND/modules/user/api/endpointsconnect"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func AuthenticationInterceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {

			anonymousPaths := map[string]interface{}{

				"/oauth/google/callback": nil,

				user_endpointsconnect.UserServiceCreateUserViaEmailHandlerProcedure: nil,
				user_endpointsconnect.UserServiceLoginViaPasswordHandlerProcedure:   nil,
				user_endpointsconnect.UserServiceForgotPasswordHandlerProcedure:     nil,
				user_endpointsconnect.UserServiceResetPasswordHandlerProcedure:      nil,
				user_endpointsconnect.UserServiceLoginViaOAuthHandlerProcedure:      nil,

				organization_endpointsconnect.OrganizationServiceGetOrganizationHandlerProcedure: nil,
				// organization_endpointsconnect.OrganizationServiceGetSalesUnitsHandlerProcedure:   nil,

				order_endpointsconnect.OrderServiceConfirmBankTransferHandlerProcedure:             nil,
				order_endpointsconnect.OrderServiceGetEggOrderTotalHandlerProcedure:                nil,
				order_endpointsconnect.OrderServiceGetFreightOrderTotalHandlerProcedure:            nil,
				order_endpointsconnect.OrderServiceGetOrderHandlerProcedure:                        nil,
				order_endpointsconnect.OrderServiceGetPaypalOrderHandlerProcedure:                  nil,
				order_endpointsconnect.OrderServiceGetRentalVehicleOrderTotalHandlerProcedure:      nil,
				order_endpointsconnect.OrderServiceCreatePotableWaterOrderHandlerProcedure:         nil,
				order_endpointsconnect.OrderServiceCreateIrrigationWaterOrderHandlerProcedure:      nil,
				order_endpointsconnect.OrderServiceGetOrderTypesAvailableInCountryHandlerProcedure: nil,
				order_endpointsconnect.OrderServiceGetDiscountCodesHandlerProcedure:                nil,

				logistic_endpointsconnect.LogisticServiceCompleteLogisticHandlerProcedure:        nil,
				logistic_endpointsconnect.LogisticServiceGetLogisticQueueSizeHandlerProcedure:    nil,
				logistic_endpointsconnect.LogisticServiceGetLogisticScheduleHandlerProcedure:     nil,
				logistic_endpointsconnect.LogisticServiceGetRentalVehiclesNearbyHandlerProcedure: nil,
				logistic_endpointsconnect.LogisticServiceGetBulkGoodTrucksNearbyHandlerProcedure: nil,
				logistic_endpointsconnect.LogisticServiceGetFreightTrucksNearbyHandlerProcedure:  nil,
				// logistic_endpointsconnect.LogisticServiceGetTrucksHandlerProcedure:                    nil,
				logistic_endpointsconnect.LogisticServiceReverseGeocodingHandlerProcedure:             nil,
				logistic_endpointsconnect.LogisticServiceGeocodingHandlerProcedure:                    nil,
				logistic_endpointsconnect.LogisticServiceValidateAddressesHandlerProcedure:            nil,
				logistic_endpointsconnect.LogisticServiceGetFreightTruckClassesNearbyHandlerProcedure: nil,

				payment_endpointsconnect.PaymentServicePaypalCaptureOrderHandlerProcedure:      nil,
				payment_endpointsconnect.PaymentServiceBitcoinPaymentHandlerProcedure:          nil,
				payment_endpointsconnect.PaymentServiceBitcoinLightningPaymentHandlerProcedure: nil,
			}

			_, ok := anonymousPaths[req.Spec().Procedure]
			if ok {
				return next(ctx, req)
			}

			authHeaderValue := req.Header().Get("Authorization")
			if authHeaderValue == "" {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("authorization header is missing"),
				)
			}

			claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("invalid authorization token provided"),
				)
			}

			if claims.Google != nil && time.Now().UnixMilli() >= claims.Google.Expires {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("authorization token is expired"),
				)
			}

			if claims.Google == nil && time.Now().UnixMilli() >= claims.ExpiresAt.UnixMilli() {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("authorization token is expired"),
				)
			}

			return next(ctx, req)
		}
	}
}
