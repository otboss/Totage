package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	db_shared "waterordersystem/GOBACKEND/shared/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func setup_TestLoginViaPasswordHandler(t *testing.T) string {
	t.Helper() // Marks this as a helper for better error reporting
	dbState := "initialized"

	err := db_user.Query.DeleteEmailVerificationCodeByEmail(context.Background(), TestStubs.TestUserEmail)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	client := &http.Client{}
	requestBody := endpoints.CreateUserViaEmailCommand{
		Email: TestStubs.TestUserEmail,
	}

	requestBodyJson, err := json.Marshal(&requestBody)
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/CreateUserViaEmailHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("failed to create user via email: %v", err)
	}
	defer resp.Body.Close()

	var responseBody endpoints.CreateUserViaEmailResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	u, err := uuid.Parse(responseBody.VerificationSessionToken)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	verificationCodefromDb, err := db_shared.Query.GetVerificationCode(context.Background(), pgtype.UUID{
		Bytes: u,
		Valid: true,
	})
	if err != nil {
		t.Errorf("failed to get verification code from db: %v", err)
	}

	TestStubs.VerificationCode = verificationCodefromDb.VerificationCode
	TestStubs.SessionToken = responseBody.VerificationSessionToken

	info := &TestStubs.UserInfo
	info.VerificationSession = &endpoints.CreateUserViaEmailCommand_UserInfo_VerificationSession{
		SessionToken:     TestStubs.SessionToken,
		VerificationCode: TestStubs.VerificationCode,
	}
	requestBody.UserInfo = info

	requestBodyJson, err = json.Marshal(&requestBody)
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/CreateUserViaEmailHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// var responseBody endpoints.CreateUserViaEmailResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK || responseBody.JwtToken == nil {
		t.Errorf("failed to log in, incorrect status code or jwt token was not returned")
	}

	claims, err := utils.JWT.VerifyToken(*responseBody.JwtToken, utils.Env.ENC_KEY)
	if err != nil {
		t.Errorf("failed to parse Jwt token: %v", err)
	}

	if responseBody.VerificationCode == nil {
		t.Error("verification code for password reset is null")
	}

	TestStubs.JwtToken = *responseBody.JwtToken
	TestStubs.TestUserId = claims.UserID

	//
	//
	//
	//

	requestBodyResetPassword := endpoints.ResetPasswordCommand{
		SessionToken:     responseBody.VerificationSessionToken,
		VerificationCode: *responseBody.VerificationCode,
		Password:         TestStubs.TestUserPassword,
	}

	requestBodyJson, err = json.Marshal(&requestBodyResetPassword)
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	req, _ = http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/ResetPasswordHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("error with forgot password request")
	}

	return dbState
}

func TestLoginViaPasswordHandler_LoginWithPassword(t *testing.T) {
	state := setup_TestLoginViaPasswordHandler(t)

	if state != "initialized" {
		t.Errorf("expected initialized, got %s", state)
	}

	// TODO: Generate random string as new password and update

	client := &http.Client{}
	requestBody := endpoints.LoginViaPasswordCommand{
		Email:    TestStubs.TestUserEmail,
		Password: TestStubs.TestUserPassword,
	}

	requestBodyJson, err := json.Marshal(&requestBody)
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/LoginViaPasswordHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("failed to create user via email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("error changing user password")
	}
}
