package application

import (
	"context"
	"errors"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
)

func (s *Handlers) GetOrganizationHandler(
	ctx context.Context,
	req *connect.Request[endpoints.GetOrganizationQuery],
) (*connect.Response[endpoints.GetOrganizationResponse], error) {

	// authHeaderValue := req.Header().Get("Authorization")
	// if authHeaderValue == "" {
	// 	return nil, errors.New("authorization header is missing")
	// }

	// claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	// if err != nil {
	// 	return nil, errors.New("invalid authorization token provided")
	// }

	organization, err := db_organization.Query.GetOrganization(ctx, req.Msg.OrganizationId)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	country, exists := common.CountryCodeAlpha3_value[string(organization.Country)]
	if exists == false {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	currency, exists := common.Currency_value[string(organization.Currency)]
	if exists == false {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	countryCode := common.CountryCodeAlpha3(country)
	currencyCode := common.Currency(currency)

	maxOrderQueueSize := int64(organization.MaxOrderQueueSize.Int16)
	var bankAccountDetails []*endpoints.CreateOrganizationCommand_BankAccountConfig

	for _, bankAccountDetail := range bankAccountDetails {
		bankAccountDetails = append(bankAccountDetails, &endpoints.CreateOrganizationCommand_BankAccountConfig{
			Bank:             bankAccountDetail.Bank,
			Branch:           bankAccountDetail.Branch,
			AccountName:      bankAccountDetail.AccountName,
			AccountNumber:    bankAccountDetail.AccountNumber,
			AccountType:      bankAccountDetail.AccountType,
			AccountOwnership: bankAccountDetail.AccountOwnership,
		})
	}

	var decryptedRestSdkClientID *string = nil
	paypalCredentials, err := db_organization.Query.GetPayPalCredentialByOrganizationID(ctx, req.Msg.OrganizationId)
	if err != nil && err != pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}
	if err != pgx.ErrNoRows {
		val, err := utils.AES.Decrypt(utils.Env.ENC_KEY, paypalCredentials.RestSdkClientIDEncrypted.String)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}
		decryptedRestSdkClientID = &val
	}

	return connect.NewResponse(&endpoints.GetOrganizationResponse{
		Organization: &endpoints.CreateOrganizationCommand{
			OrganizationId:     &organization.ID,
			OrganizationName:   organization.OrganizationName,
			Address:            organization.Address,
			Phone:              organization.Phone,
			Country:            countryCode,
			Currency:           currencyCode,
			DutyFreeKm:         organization.DutyFreeKm,
			Email:              organization.Email,
			MaxOrderQueueSize:  &maxOrderQueueSize,
			SubmitForApproval:  &organization.IsSubmittedForApproval,
			BankAccountDetails: bankAccountDetails,
			PaypalMerchantDetails: &endpoints.PaypalMerchantDetails{
				SubscriptionId: paypalCredentials.SubscriptionID,
			},
		},
		PaypalClientId: decryptedRestSdkClientID,
		// IsApproved:     organization.IsApproved,
		Rating: int64(organization.Rating),
	}), nil
}
