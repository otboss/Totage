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
	"waterordersystem/GOBACKEND/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

func setup_TestForgotPasswordHandler(t *testing.T) string {
	t.Helper()
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

	// Delete existing verification session
	err = db_user.Query.DeleteVerificationSessionByEmail(context.Background(), pgtype.Text{
		String: TestStubs.TestUserEmail,
		Valid:  true,
	})
	if err != nil {
		t.Errorf("error occurred: %v", err)
	}

	return dbState
}

func TestForgotPasswordHandler_SendForgotPasswordEmail(t *testing.T) {
	state := setup_TestForgotPasswordHandler(t)

	if state != "initialized" {
		t.Errorf("expected initialized, got %s", state)
	}

	requestBodyForgotPassword := endpoints.ForgotPasswordCommand{
		Email: TestStubs.TestUserEmail,
	}

	client := &http.Client{}
	requestBodyJson, err := json.Marshal(&requestBodyForgotPassword)
	if err != nil {
		t.Errorf("Parse error: %v", err)
	}

	for x := 0; x < 2; x++ {
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v1.UserService/ForgotPasswordHandler", utils.Env.BACKEND_ORIGIN), bytes.NewBuffer(requestBodyJson))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if x == 0 {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("error with initial forgot password request")
			}
			continue
		}

		if resp.StatusCode == http.StatusOK {
			t.Errorf("error: forgot password requests in short succession should be blocked")
		}
	}
}
