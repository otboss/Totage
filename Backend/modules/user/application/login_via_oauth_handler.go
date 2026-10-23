package application

import (
	"context"
	"errors"
	"strconv"
	"time"
	endpoints "waterordersystem/GOBACKEND/modules/user/api"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
	"waterordersystem/GOBACKEND/shared/infrastructure/clients"
	"waterordersystem/GOBACKEND/utils"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const callbackRoute string = "/oauth/google/callback"

// func init() {
// 	googleOauthConfig := &oauth2.Config{
// 		RedirectURL:  fmt.Sprintf("%s%s", utils.Env.BACKEND_ORIGIN, callbackRoute),
// 		ClientID:     utils.Env.OAUTH_GOOGLE_CLIENT_ID,
// 		ClientSecret: utils.Env.OAUTH_GOOGLE_CLIENT_SECRET,
// 		Scopes:       []string{"https://www.googleapis.com", "https://www.googleapis.com"},
// 		Endpoint:     google.Endpoint,
// 	}

// 	// Plain HTTP route for OAuth redirect
// 	utils.Mux.HandleFunc(fmt.Sprintf("/oauth/google/callback"), func(w http.ResponseWriter, r *http.Request) {
// 		// 1. Verify 'state' parameter to prevent CSRF
// 		code := r.FormValue("code")
// 		ctx := context.Background()

// 		// 2. Exchange authorization code for token
// 		token, err := googleOauthConfig.Exchange(ctx, code)
// 		if err != nil {
// 			http.Error(w, "Code exchange failed", http.StatusInternalServerError)
// 			return
// 		}

// 		// 3. (Optional) Use token to get user info
// 		client := googleOauthConfig.Client(ctx, token)
// 		resp, err := client.Get("https://www.googleapis.com")
// 		if err != nil {
// 			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
// 			return
// 		}

// 		var googleUser struct {
// 			Sub           string `json:"sub"`
// 			Name          string `json:"name"`
// 			GivenName     string `json:"given_name"`
// 			FamilyName    string `json:"family_name"`
// 			Picture       string `json:"picture"`
// 			Email         string `json:"email"`
// 			EmailVerified bool   `json:"email_verified"`
// 		}
// 		if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
// 			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
// 			return
// 		}

// 		user, err := db_user.Query.GetUserByEmail(ctx, googleUser.Email)
// 		if err == pgx.ErrNoRows {
// 			encryptedRefreshToken, err := utils.AES.Encrypt(utils.Env.ENC_KEY, token.RefreshToken)
// 			if err != nil {
// 				http.Error(w, "internal server error", http.StatusInternalServerError)
// 			}

// 			userId, err := db_user.Query.CreateUser(ctx, db_user.CreateUserParams{
// 				FirstName: googleUser.GivenName,
// 				LastName:  googleUser.FamilyName,
// 				Email: pgtype.Text{
// 					String: googleUser.Email,
// 					Valid:  true,
// 				},
// 				GoogleUserID: pgtype.Text{
// 					String: googleUser.Sub,
// 					Valid:  true,
// 				},
// 				OauthRefreshTokenEncrypted: pgtype.Text{
// 					String: encryptedRefreshToken,
// 					Valid:  true,
// 				},
// 				LastIpAddress: r.RemoteAddr,
// 			})
// 			if err != nil {
// 				http.Error(w, "internal server error", http.StatusInternalServerError)
// 			}

// 			user.ID = userId
// 		}
// 		if err != nil && err != pgx.ErrNoRows {
// 			http.Error(w, "internal server error", http.StatusInternalServerError)
// 		}

// 		jwtToken, err := utils.JWT.CreateToken(user.ID, utils.Env.ENC_KEY)
// 		if err != nil {
// 			http.Error(w, "internal server error", http.StatusInternalServerError)
// 		}

// 		http.Redirect(w, r, fmt.Sprintf("%s#access_token=%s", utils.Env.FRONTEND_ORIGIN, jwtToken), http.StatusTemporaryRedirect)
// 	})
// }

func (s *Handlers) LoginViaOAuthHandler(
	ctx context.Context,
	req *connect.Request[endpoints.LoginViaOAuthCommand],
) (*connect.Response[endpoints.LoginViaOAuthResponse], error) {

	oauthInfo, err := clients.GoogleCloudApi.OAuthTokenInfo(req.Msg.AccessToken)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if oauthInfo.EmailVerified != "true" {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("email unverified"),
		)
	}

	expiresIn, err := strconv.Atoi(oauthInfo.ExpiresIn)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if expiresIn <= 0 {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("invalid oauth token provided"),
		)
	}

	_, err = db_user.Query.GetUserByEmail(ctx, oauthInfo.Email)
	if err == pgx.ErrNoRows {
		userInfo, err := clients.GoogleCloudApi.OAuthUserInfo(req.Msg.AccessToken)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		_, err = db_user.Query.CreateUser(ctx, db_user.CreateUserParams{
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
			Email: pgtype.Text{
				String: userInfo.Email,
				Valid:  true,
			},
			GoogleUserID: pgtype.Text{
				String: userInfo.ID,
				Valid:  true,
			},
		})

		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

	} else if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	user, err := db_user.Query.GetUserByEmail(ctx, oauthInfo.Email)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	jwtToken, err := utils.JWT.CreateToken(user.ID, &struct {
		AccessToken string "json:\"google_access_token\""
		Expires     int64  "json:\"expires\""
	}{
		AccessToken: req.Msg.AccessToken,
		Expires:     time.Now().UnixMilli() + int64(expiresIn*1000),
	}, utils.Env.ENC_KEY)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.LoginViaOAuthResponse{
		JwtToken: jwtToken,
	}), nil
}
