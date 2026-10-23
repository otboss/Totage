package application

import "testing"

func setup_TestLoginViaOAuthHandler(t *testing.T) string {
	dbState := "initialized"

	// TODO: Create new user

	return dbState
}

func TestLoginViaOAuthHandler_LoginWithOAuth(t *testing.T) {
	state := setup_TestForgotPasswordHandler(t)

	if state != "initialized" {
		t.Errorf("expected initialized, got %s", state)
	}

	// TODO: Complete this function
}
