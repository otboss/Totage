package application

import (
	"fmt"
	"net/http"
	"sync"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	"waterordersystem/GOBACKEND/modules/user/api/endpointsconnect"
	"waterordersystem/GOBACKEND/shared/infrastructure/middleware"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgtype"
)

func init() {
	handlers := &Handlers{}

	var path string
	var handler http.Handler

	path, handler = endpointsconnect.NewUserServiceHandler(
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

type totalDistanceInKm struct {
	mu    sync.Mutex
	total int32
}

func (t *totalDistanceInKm) Add(distance int32) {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.total += distance
}

type testStubs struct {
	SessionToken        string
	VerificationCodesId pgtype.UUID
	VerificationCode    int32
	TestUserId          int64
	TestUserEmail       string
	TestUserPassword    string
	UserInfo            endpoints.CreateUserViaEmailCommand_UserInfo
	JwtToken            string
}

var TestStubs *testStubs = &testStubs{
	SessionToken:     "",
	VerificationCode: 0,
	TestUserEmail:    "johndoe@gmail.com",
	TestUserPassword: "Password@123",
	UserInfo: endpoints.CreateUserViaEmailCommand_UserInfo{
		FirstName:      "John",
		LastName:       "Doe",
		Phone:          "18765555555",
		RecaptchaEvent: nil,
		// Password:       "Password@123",
		// VerificationSession: &verificationSession,
	},
	JwtToken: "",
}
