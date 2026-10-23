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

func setup_TestCreateUserViaEmailHandler(t *testing.T) string {
	t.Helper() // Marks this as a helper for better error reporting
	dbState := "initialized"

	// t.Cleanup acts as your "AfterEach"
	// It runs automatically after the test (and any subtests) finishes.
	// t.Cleanup(func() {
	// 	dbState = "cleaned"
	// })

	err := db_user.Query.DeleteEmailVerificationCodeByEmail(context.Background(), TestStubs.TestUserEmail)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// if state != "initialized" {
	// 	t.Errorf("expected initialized, got %s", state)
	// }

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

	return dbState
}

func TestCreateUserViaEmailHandler_VerifySessionWithIncorrectVerificationCode(t *testing.T) {
	state := setup_TestCreateUserViaEmailHandler(t)

	if state != "initialized" {
		t.Errorf("expected initialized, got %s", state)
	}

	info := &TestStubs.UserInfo
	info.VerificationSession = &endpoints.CreateUserViaEmailCommand_UserInfo_VerificationSession{
		SessionToken:     TestStubs.SessionToken,
		VerificationCode: 0,
	}

	client := &http.Client{}
	requestBody := endpoints.CreateUserViaEmailCommand{
		Email:    TestStubs.TestUserEmail,
		UserInfo: info,
	}

	requestBodyJson, err := json.Marshal(&requestBody)
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/CreateUserViaEmailHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("user session incorrectly created")
	}
}

func TestCreateUserViaEmailHandler_VerifySessionWithCorrectVerificationCode(t *testing.T) {
	state := setup_TestCreateUserViaEmailHandler(t)

	if state != "initialized" {
		t.Errorf("expected initialized, got %s", state)
	}

	info := &TestStubs.UserInfo
	info.VerificationSession = &endpoints.CreateUserViaEmailCommand_UserInfo_VerificationSession{
		SessionToken:     TestStubs.SessionToken,
		VerificationCode: TestStubs.VerificationCode,
	}

	client := &http.Client{}
	requestBody := endpoints.CreateUserViaEmailCommand{
		Email:    TestStubs.TestUserEmail,
		UserInfo: info,
	}

	requestBodyJson, err := json.Marshal(&requestBody)
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/CreateUserViaEmailHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var responseBody endpoints.CreateUserViaEmailResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK || responseBody.JwtToken == nil {
		t.Errorf("failed to log in, incorrect status code or jwt token was not returned")
	}
}

func TestCreateUserViaEmailHandler_VerifySessionWithIncorrectVerificationCodeTooManyTimes(t *testing.T) {
	state := setup_TestCreateUserViaEmailHandler(t)
	if state != "initialized" {
		t.Errorf("expected initialized, got %s", state)
	}

	info := &TestStubs.UserInfo
	info.VerificationSession = &endpoints.CreateUserViaEmailCommand_UserInfo_VerificationSession{
		SessionToken:     TestStubs.SessionToken,
		VerificationCode: 0,
	}

	iterations := 4

	for x := 0; x < iterations; x++ {

		client := &http.Client{}
		requestBody := endpoints.CreateUserViaEmailCommand{
			Email:    TestStubs.TestUserEmail,
			UserInfo: info,
		}

		if x == iterations-1 {
			info.VerificationSession.VerificationCode = TestStubs.VerificationCode
		}

		requestBodyJson, err := json.Marshal(&requestBody)
		if err != nil {
			t.Errorf("Parse error: %v", err)
		}

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/CreateUserViaEmailHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if x == iterations-1 {
			info.VerificationSession.VerificationCode = TestStubs.VerificationCode
			if resp.StatusCode == http.StatusOK {
				t.Errorf("user session incorrectly created")
			}
			break
		}

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("user session incorrectly created")
		}
	}

}
