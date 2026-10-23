package application

import (
	"context"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"

	"connectrpc.com/connect"
)

func (s *Handlers) VerifyAccountOAuthHandler(
	ctx context.Context,
	req *connect.Request[endpoints.VerifyAccountOAuthCommand],
) (*connect.Response[endpoints.VerifyAccountOAuthResponse], error) {

	// TODO:

	// authHeaderValue := req.Header().Get("Authorization")
	// if authHeaderValue == "" {
	// 	return nil, errors.New("authorization header is missing")
	// }

	// claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	// if err != nil {
	// 	return nil, errors.New("invalid authorization token provided")
	// }

	// if claims.UserType != utils.Business {
	// 	return nil, errors.New("action not supported")
	// }

	// accountVerificationDetails, err := db_organization.Query.GetAccountByVerificationToken(ctx, req.Msg.VerificationToken)
	// if err != nil {
	// 	return nil, connect.NewError(
	// 	connect.CodeInternal,
	// 	errors.New("internal server error"),
	// )
	// }

	// account, err := db_organization.Query.GetAccount(ctx, accountVerificationDetails.AccountID)
	// if err != nil {
	// 	return nil, connect.NewError(
	// 	connect.CodeInternal,
	// 	errors.New("internal server error"),
	// )
	// }

	// err = db_organization.Query.DeleteVerificationToken(ctx, accountVerificationDetails.ID)
	// if err != nil {
	// 	return nil, connect.NewError(
	// 	connect.CodeInternal,
	// 	errors.New("internal server error"),
	// )
	// }

	// hashedPasswordEncrypted, err := utils.AES.Encrypt(utils.Env.ENC_KEY, hashedPassword)
	// if err != nil {
	// 	return nil, connect.NewError(
	// 	connect.CodeInternal,
	// 	errors.New("internal server error"),
	// )
	// }

	// userId, err := db_organization.Query.CreateAccount(ctx, db_organization.CreateAccountParams{
	// 	OrganizationID:    account.OrganizationID,
	// 	FirstName:         req.Msg.FirstName,
	// 	LastName:          req.Msg.LastName,
	// 	Phone:             req.Msg.Phone,
	// 	Email:             req.Msg.Email,
	// 	IsAdmin:           req.Msg.IsAdmin,
	// 	IsCustomerSupport: req.Msg.IsCustomerSupport,
	// 	IsDriver:          req.Msg.IsDriver,
	// })
	// if err != nil {
	// 	return nil, connect.NewError(
	// 	connect.CodeInternal,
	// 	errors.New("internal server error"),
	// )
	// }

	// clients.GoogleCloudApi.

	return connect.NewResponse(&endpoints.VerifyAccountOAuthResponse{}), nil
}
