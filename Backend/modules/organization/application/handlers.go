package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	endpoints_organization "waterordersystem/GOBACKEND/modules/organization/api"
	"waterordersystem/GOBACKEND/modules/organization/api/endpointsconnect"
	endpoints_user "waterordersystem/GOBACKEND/modules/user/api"
	common "waterordersystem/GOBACKEND/shared/api"
	"waterordersystem/GOBACKEND/shared/infrastructure/middleware"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

func init() {
	handlers := &Handlers{}

	var path string
	var handler http.Handler

	path, handler = endpointsconnect.NewOrganizationServiceHandler(
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

type testStubs struct {
	SessionToken        string
	VerificationCodesId pgtype.UUID
	VerificationCode    int32
	TestUserId          int64
	TestUserEmail       string
	TestUserPassword    string
	AuthToken           string
	UserInfo            endpoints_user.CreateUserViaEmailCommand_UserInfo
	CreateOrganization  endpoints_organization.CreateOrganizationCommand
	JwtToken            string
}

func (t *testStubs) CreateUser() error {
	client := &http.Client{}
	requestBody := endpoints_user.CreateUserViaEmailCommand{
		Email: TestStubs.TestUserEmail,
	}

	requestBodyJson, err := json.Marshal(&requestBody)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/CreateUserViaEmailHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var responseBody endpoints_user.CreateUserViaEmailResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	//
	//
	//
	//
	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/CreateUserViaEmailHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK || responseBody.JwtToken == nil {
		return err
	}

	claims, err := utils.JWT.VerifyToken(*responseBody.JwtToken, utils.Env.ENC_KEY)
	if err != nil {
		return err
	}

	if responseBody.VerificationCode == nil {
		return err
	}

	t.JwtToken = *responseBody.JwtToken
	t.TestUserId = claims.UserID
	t.SessionToken = responseBody.VerificationSessionToken
	t.VerificationCode = *responseBody.VerificationCode

	// requestBodyResetPassword := endpoints.ResetPasswordCommand{
	// 	SessionToken:     TestStubs.SessionToken,
	// 	VerificationCode: TestStubs.VerificationCode,
	// 	Password:         TestStubs.TestUserPassword,
	// }

	// requestBodyJson, err = json.Marshal(&requestBodyResetPassword)
	// if err != nil {
	// 	return err
	// }

	// req, _ = http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/ResetPasswordHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	// req.Header.Set("Content-Type", "application/json")

	// resp, err = client.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != 200 {
	// 	return err
	// }

	return nil
}

var TestStubs *testStubs = &testStubs{
	SessionToken:     "",
	VerificationCode: 0,
	TestUserEmail:    "johndoe@gmail.com",
	TestUserPassword: "Password@123",
	AuthToken:        "",
	UserInfo: endpoints_user.CreateUserViaEmailCommand_UserInfo{
		FirstName:      "John",
		LastName:       "Doe",
		Phone:          "18765555555",
		RecaptchaEvent: nil,
		// Password:       "Password@123",
		// VerificationSession: &verificationSession,
	},
	CreateOrganization: endpoints_organization.CreateOrganizationCommand{
		OrganizationName: "Test Organization",
		Address:          "123 Test Ave.",
		Phone:            "+18765555555",
		Email:            "otsurfer7@gmail.com",
		Country:          common.CountryCodeAlpha3_JAM,
		Currency:         common.Currency_JMD,
		DutyFreeKm:       12.0,
		// FreightConfig: &endpoints_organization.CreateOrganizationCommand_FreightConfig{
		// 	MinimumJobCost: 8000,
		// 	IsEnabled:      true,
		// },
		PaypalMerchantDetails: &endpoints_organization.PaypalMerchantDetails{
			SubscriptionId: "test",
		},
		BankAccountDetails: []*endpoints_organization.CreateOrganizationCommand_BankAccountConfig{
			&endpoints_organization.CreateOrganizationCommand_BankAccountConfig{
				Bank:             "National Commercial Bank Jamaica",
				Branch:           "UWI Mona",
				AccountName:      "John Doe",
				AccountNumber:    "123456789",
				AccountType:      endpoints_organization.CreateOrganizationCommand_BankAccountConfig_SAVINGS,
				AccountOwnership: endpoints_organization.CreateOrganizationCommand_BankAccountConfig_BUSINESS,
			},
		},
		// SupportedOrder: []endpoints_organization.SupportedOrder{
		// 	endpoints_organization.SupportedOrder_FREIGHT,
		// },
		// FreightConfig: &endpoints_organization.CreateOrganizationCommand_FreightConfig{
		// 	MinimumJobCost:        8000,
		// 	MaxDutyFreeDistanceKm: 10,
		// 	DutyCostPerKm:         150,
		// },
	},
	// UpdateOrganization: endpoints_organization.UpdateOrganizationCommand{
	// JobTimeslots: []*endpoints_organization.JobTimeslot{
	// 	&endpoints_organization.JobTimeslot{
	// 		HH: 6,
	// 		MM: 00,
	// 		Weekdays: []common.Weekdays{
	// 			common.Weekdays_MONDAY,
	// 			common.Weekdays_TUESDAY,
	// 			common.Weekdays_WEDNESDAY,
	// 			common.Weekdays_THURSDAY,
	// 			common.Weekdays_FRIDAY,
	// 		},
	// 	},
	// },
	// },
	JwtToken: "",
}
