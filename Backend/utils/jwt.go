package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	// endpoints "waterordersystem/GOBACKEND/shared/api"

	"github.com/golang-jwt/jwt/v5"
)

type jsonWebToken struct {
}

// type UserType int

// const (
// 	Customer UserType = iota
// 	Business
// )

// type Role string
// const(
// 	Administrator Role = "ADMINISTRATOR"
// 	Driver Role = "DRIVER"

// )
// type OrganizationId int64

type CustomClaims struct {
	UserID int64 `json:"user_id"`
	Google *struct {
		AccessToken string `json:"google_access_token"`
		Expires     int64  `json:"expires"`
	}
	// UserType UserType `json:"user_type"`
	jwt.RegisteredClaims
	// OrganizationAccountRoles map[OrganizationId]endpoints.AccountRole
}

// func (c UserType) String() string {
// 	switch c {
// 	case Customer:
// 		return "CUSTOMER"
// 	case Business:
// 		return "BUSINESS"
// 	default:
// 		return fmt.Sprintf("UserType(%d)", c)
// 	}
// }

// createToken generates a new JWT token with custom claims
func (j *jsonWebToken) CreateToken(
	userID int64,
	googleAccessToken *struct {
		AccessToken string `json:"google_access_token"`
		Expires     int64  `json:"expires"`
	},
	signingKey string,
) (string, error) {

	claims := CustomClaims{
		UserID: userID,
		Google: googleAccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my-app",
			Subject:   "user-auth",
		},
	}

	// Create the token with the signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key to get the complete token string
	signedTokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedTokenString, nil
}

// verifyToken parses and validates a JWT token string
func (j *jsonWebToken) VerifyToken(tokenString, signingKey string) (CustomClaims, error) {
	if strings.Index(
		strings.ToLower(tokenString),
		"bearer ",
	) == 0 {
		tokenString = tokenString[7:]
	}

	var customClaims CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return the secret key for verification
		return []byte(signingKey), nil
	})
	if err != nil {
		return customClaims, fmt.Errorf("token parsing or validation failed: %w", err)
	}

	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return customClaims, err
	}

	if expirationTime.Before(time.Now()) == true {
		return customClaims, errors.New("token expired please log in")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return *claims, nil
	} else {
		return customClaims, errors.New("invalid token claims or token is invalid")
	}
}

var JWT jsonWebToken = jsonWebToken{}
