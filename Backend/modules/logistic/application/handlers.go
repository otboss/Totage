package application

import (
	"fmt"
	"net/http"
	"waterordersystem/GOBACKEND/modules/logistic/api/endpointsconnect"
	"waterordersystem/GOBACKEND/shared/infrastructure/middleware"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
)

func init() {
	handlers := &Handlers{}

	var path string
	var handler http.Handler

	path, handler = endpointsconnect.NewLogisticServiceHandler(
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

type helperFunctions struct{}

var HelperFunctions = helperFunctions{}
