package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	endpoints_user "waterordersystem/GOBACKEND/modules/user/api"
	db_shared "waterordersystem/GOBACKEND/shared/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func setup_TestCreateOrganizationHandler(t *testing.T) string {
	t.Helper()
	dbState := "initialized"

	err := TestStubs.CreateUser()
	if err != nil {
		t.Errorf("Error while creating user: %v", err)
	}

	err = db_organization.Query.DeleteOrganizations(context.Background())
	if err != nil {
		t.Errorf("Error clearing organizations: %v", err)
	}

	// TODO: Create user and sET TestStubs.Auth Token
	client := &http.Client{}
	requestBody := endpoints_user.CreateUserViaEmailCommand{
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

	var responseBody endpoints_user.CreateUserViaEmailResponse
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
	info.VerificationSession = &endpoints_user.CreateUserViaEmailCommand_UserInfo_VerificationSession{
		SessionToken:     TestStubs.SessionToken,
		VerificationCode: TestStubs.VerificationCode,
	}

	requestBody = endpoints_user.CreateUserViaEmailCommand{
		Email:    TestStubs.TestUserEmail,
		UserInfo: info,
	}

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

	TestStubs.AuthToken = *responseBody.JwtToken

	// User created successfully and auth token set
	return dbState
}

func TestCreateOrganizationHandler_CreateOrganization(t *testing.T) {
	state := setup_TestCreateOrganizationHandler(t)

	if state != "initialized" {
		t.Errorf("expected initialized, got %s", state)
	}

	client := &http.Client{}

	requestBodyJson, err := json.Marshal(&TestStubs.CreateOrganization)
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1.OrganizationService/CreateOrganizationHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
	// TODO: Add authorization header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", TestStubs.AuthToken))

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("failed to create user via email: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Convert bytes to plain text string
	bodyString := string(bodyBytes)
	t.Errorf("%s", bodyString)
	fmt.Println(bodyString)

	var responseBody endpoints.CreateOrganizationResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Errorf("failed to decode response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("error while create creating organization")
	}
}
