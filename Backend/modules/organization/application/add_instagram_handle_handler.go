package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	endpoints "waterordersystem/GOBACKEND/modules/organization/api"
	db_organization "waterordersystem/GOBACKEND/modules/organization/infrastructure"
	"waterordersystem/GOBACKEND/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"

	"connectrpc.com/connect"
)

func (s *Handlers) AddInstagramHandleHandler(
	ctx context.Context,
	req *connect.Request[endpoints.AddInstagramHandleCommand],
) (*connect.Response[endpoints.AddInstagramHandleResponse], error) {
	authHeaderValue := req.Header().Get("Authorization")
	if authHeaderValue == "" {
		return nil, errors.New("authorization header is missing")
	}

	claims, err := utils.JWT.VerifyToken(authHeaderValue, utils.Env.ENC_KEY)
	if err != nil {
		return nil, errors.New("invalid authorization token provided")
	}

	account, err := db_organization.Query.GetAccountByUserIdAndOrganizationId(ctx, db_organization.GetAccountByUserIdAndOrganizationIdParams{
		UserID:         claims.UserID,
		OrganizationID: req.Msg.OrganizationId,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	if account.IsAdmin == false {
		return nil, errors.New("not permitted")
	}

	instagramHandleVerification, err := db_organization.Query.GetInstagramHandleVerificationByOrganizationId(ctx, account.OrganizationID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	isCurrentTimeBeforeExpiration := time.Now().Before(instagramHandleVerification.CreatedAt.Time.Add(time.Duration(utils.Env.INSTAGRAM_HANDLE_VERIFICATION_TIME_MS) * time.Millisecond))
	expirationTime := instagramHandleVerification.CreatedAt.Time.Add(time.Duration(utils.Env.INSTAGRAM_HANDLE_VERIFICATION_TIME_MS) * time.Millisecond)
	expirationTimestampProto := timestamppb.New(expirationTime)

	if err == pgx.ErrNoRows || isCurrentTimeBeforeExpiration == false {
		// FORCE START NEW VERIFICATION SESSION
		instagramHandleVerification, err := db_organization.Query.CreateInstagramHandleVerification(ctx, account.OrganizationID)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		expirationTime := time.Now().Add(time.Duration(utils.Env.INSTAGRAM_HANDLE_VERIFICATION_TIME_MS) * time.Millisecond)
		expirationTimeProto := timestamppb.New(expirationTime)

		return connect.NewResponse(&endpoints.AddInstagramHandleResponse{
			VerificationCode:         instagramHandleVerification.InstagramHandleVerificationCode,
			IsVerificationSuccessful: false,
			FailedAttempts:           0,
			VerificationEndTime:      expirationTimeProto,
		}), nil
	}

	if instagramHandleVerification.FailedAttempts >= 3 && isCurrentTimeBeforeExpiration == true {
		return nil, connect.NewError(
			connect.CodeResourceExhausted,
			errors.New("too many verification attempts, please try again later"),
		)
	}

	if req.Msg.InstagramHandle == nil && isCurrentTimeBeforeExpiration == true {
		return connect.NewResponse(&endpoints.AddInstagramHandleResponse{
			VerificationCode:         instagramHandleVerification.InstagramHandleVerificationCode,
			IsVerificationSuccessful: false,
			FailedAttempts:           int64(instagramHandleVerification.FailedAttempts),
			VerificationEndTime:      expirationTimestampProto,
		}), nil
	}

	bio, avatarUrl, err := getInstagramBioAndAvatar(*req.Msg.InstagramHandle)
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	isVerificationSuccessful := strings.Contains(bio, instagramHandleVerification.InstagramHandleVerificationCode)

	if isVerificationSuccessful == false {
		err = db_organization.Query.IncrementInstagramHandleVerificationFailedAttempt(ctx, instagramHandleVerification.ID)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeInternal,
				errors.New("internal server error"),
			)
		}

		return connect.NewResponse(&endpoints.AddInstagramHandleResponse{
			VerificationCode:         instagramHandleVerification.InstagramHandleVerificationCode,
			IsVerificationSuccessful: false,
			FailedAttempts:           int64(instagramHandleVerification.FailedAttempts + 1),
			VerificationEndTime:      expirationTimestampProto,
		}), nil
	}

	err = db_organization.Query.UpdateInstagramHandleAndAvatarUrl(ctx, db_organization.UpdateInstagramHandleAndAvatarUrlParams{
		InstagramHandle: *req.Msg.InstagramHandle,
		AvatarUrl: pgtype.Text{
			String: avatarUrl,
			Valid:  true,
		},
		OrganizationID: req.Msg.OrganizationId,
	})
	if err != nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("internal server error"),
		)
	}

	return connect.NewResponse(&endpoints.AddInstagramHandleResponse{
		VerificationCode:         instagramHandleVerification.InstagramHandleVerificationCode,
		IsVerificationSuccessful: isVerificationSuccessful,
		FailedAttempts:           int64(instagramHandleVerification.FailedAttempts),
		VerificationEndTime:      expirationTimestampProto,
	}), nil
}

//
//
//
//
//
//

func getInstagramBioAndAvatar(username string) (string, string, error) {

	// Fixed the URL formatting
	url := fmt.Sprintf("https://www.instagram.com/%s", username)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	// 1. Extract the Bio from the description meta tag
	bioContent, bioExists := doc.Find("meta[name='description']").Attr("content")
	if !bioExists {
		return "", "", fmt.Errorf("could not find meta description tag")
	}

	const marker = "on Instagram: "
	parts := strings.Split(bioContent, marker)
	bio := ""
	if len(parts) >= 2 {
		bio = strings.TrimSpace(parts[1])
	}

	// 2. Extract the Avatar URL from the og:image meta tag
	avatarURL, avatarExists := doc.Find("meta[property='og:image']").Attr("content")
	if !avatarExists {
		// Fallback: some pages might use name instead of property
		avatarURL, _ = doc.Find("meta[name='og:image']").Attr("content")
	}

	return bio, avatarURL, nil
}
