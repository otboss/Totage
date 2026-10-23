package clients

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	currencyCodes "waterordersystem/GOBACKEND/shared/domain/enums"
	"waterordersystem/GOBACKEND/utils"
)

type paypalApi struct{}

// GenerateAccessToken fetches an access token from the PayPal API.
func generateAccessToken(paypalRestSdkClientId string, paypalRestSdkSecret string) (string, error) {
	// 1. Encode client ID and secret in Base64 for Basic Auth header.
	authString := fmt.Sprintf("%s:%s", paypalRestSdkClientId, paypalRestSdkSecret)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(authString))
	// The Authorization header value will be "Basic <base64_encoded_string>"
	authHeader := fmt.Sprintf("Basic %s", authEncoded)

	// 2. Prepare the request body as form data
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	requestBody := strings.NewReader(data.Encode())

	// 3. Create the HTTP request
	requestURL := fmt.Sprintf("%s/v1/oauth2/token", utils.Env.PAYPAL_API_BASE_URL)
	req, err := http.NewRequest("POST", requestURL, requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 4. Set headers
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // Important for form data
	req.Header.Set("Accept", "application/json")

	// 5. Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	// 6. Decode the JSON response
	var responseData accessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// 7. Return the access token
	return responseData.AccessToken, nil
}

func (p paypalApi) CreateOrder(paypalRestSdkClientId string, paypalRestSdkSecret string, totalInUSD float64) (PayPalCreateOrderResponse, error) {
	response := PayPalCreateOrderResponse{}
	accessToken, err := generateAccessToken(paypalRestSdkClientId, paypalRestSdkSecret)
	if err != nil {
		return response, err
	}

	requestURL := fmt.Sprintf("%s/v2/checkout/orders", utils.Env.PAYPAL_API_BASE_URL)

	request := paypalCreateOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []purchaseUnit{
			{
				amount{
					CurrencyCode: currencyCodes.JMD,
					Value:        fmt.Sprintf("%.2f", totalInUSD),
				},
			},
		},
	}

	requestJson, err := json.Marshal(request)
	if err != nil {
		return response, err
	}

	requestBody := bytes.NewReader(requestJson)

	req, err := http.NewRequest("POST", requestURL, requestBody)
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return response, fmt.Errorf("API request failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

func (p paypalApi) CaptureOrder(paypalRestSdkClientId string, paypalRestSdkSecret string, paypalOrderId string) (PayPalCreateOrderResponse, error) {
	var response PayPalCreateOrderResponse

	accessToken, err := generateAccessToken(paypalRestSdkClientId, paypalRestSdkSecret)
	if err != nil {
		return response, err
	}

	requestURL := fmt.Sprintf("%s/v2/checkout/orders/$%s/capture", utils.Env.PAYPAL_API_BASE_URL, paypalOrderId)

	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return response, fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return response, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return response, nil
}

func (p paypalApi) ShowSubscriptionDetails(paypalRestSdkClientId string, paypalRestSdkSecret string, subscriptionId string) (PaypalShowSubscriptionDetailsResponse, error) {
	response := PaypalShowSubscriptionDetailsResponse{}
	accessToken, err := generateAccessToken(paypalRestSdkClientId, paypalRestSdkSecret)
	if err != nil {
		return response, err
	}

	requestURL := fmt.Sprintf("%s/v1/billing/subscriptions/%s", utils.Env.PAYPAL_API_BASE_URL, subscriptionId)

	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		return response, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return response, fmt.Errorf("API request failed with status %s: %s", resp.Status, string(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}

var PaypalApi = paypalApi{}

type accessTokenResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
	Nonce       string `json:"nonce"`
}

type amount struct {
	CurrencyCode currencyCodes.Currency `json:"currency_code"` // e.g., "USD"
	Value        string                 `json:"value"`         // Store value as string to match JS .toFixed(2) output
}

// PurchaseUnit defines a single unit in the purchase.
type purchaseUnit struct {
	Amount amount `json:"amount"`
}

// OrderRequest defines the top-level structure for the PayPal order creation request.
type paypalCreateOrderRequest struct {
	Intent        string         `json:"intent"` // e.g., "CAPTURE"
	PurchaseUnits []purchaseUnit `json:"purchase_units"`
}

type PayPalCreateOrderResponse struct {
	ID            string         `json:"id"`
	Status        string         `json:"status"`
	Links         []Link         `json:"links"`
	Details       []Detail       `json:"details"`
	DebugID       *string        `json:"debug_id,omitempty"` // Use pointer/omitempty for optional fields
	PurchaseUnits []PurchaseUnit `json:"purchase_units"`
}

// Detail corresponds to the inner details array.
type Detail struct {
	Issue       string `json:"issue"`
	Description string `json:"description"`
}

// PurchaseUnit corresponds to the purchase_units array elements.
type PurchaseUnit struct {
	ReferenceID string   `json:"reference_id"`
	Shipping    Shipping `json:"shipping"`
	Payments    Payments `json:"payments"`
}

// Shipping contains address details.
type Shipping struct {
	Address Address `json:"address"`
}

// Address details.
type Address struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	AdminArea2   string `json:"admin_area_2"`
	AdminArea1   string `json:"admin_area_1"`
	PostalCode   string `json:"postal_code"`
	CountryCode  string `json:"country_code"`
}

// Payments aggregates captures and authorizations.
type Payments struct {
	Captures       []Capture     `json:"captures"`
	Authorizations []interface{} `json:"authorizations"` // Use []interface{} if structure is unknown
}

// Capture represents a single payment capture.
type Capture struct {
	ID                        string                    `json:"id"`
	Status                    string                    `json:"status"`
	Amount                    Money                     `json:"amount"`
	SellerProtection          SellerProtection          `json:"seller_protection"`
	FinalCapture              bool                      `json:"final_capture"`
	DisbursementMode          string                    `json:"disbursement_mode"`
	SellerReceivableBreakdown SellerReceivableBreakdown `json:"seller_receivable_breakdown"`
	CreateTime                time.Time                 `json:"create_time"` // Use time.Time for standard date parsing
	UpdateTime                time.Time                 `json:"update_time"`
	Links                     []Link                    `json:"links"`
}

// Money represents currency and value. (Reused from your previous request)
type Money struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

// SellerProtection details.
type SellerProtection struct {
	Status            string   `json:"status"`
	DisputeCategories []string `json:"dispute_categories"`
}

// SellerReceivableBreakdown details.
type SellerReceivableBreakdown struct {
	GrossAmount Money `json:"gross_amount"`
	PaypalFee   Money `json:"paypal_fee"`
	NetAmount   Money `json:"net_amount"`
}

// Link represents HATEOAS links provided by PayPal.
type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type PaypalSubscriptionStatus string

// type paypalSubscriptionStatuses struct {
// 	Active   string
// 	Created  string
// 	InActive string
// }

const (
	// By assigning the type to the first constant,
	// subsequent constants in the block automatically inherit the type and increment
	PaypalSubscriptionStatusActive   PaypalSubscriptionStatus = "ACTIVE"
	PaypalSubscriptionStatusCreated  PaypalSubscriptionStatus = "CREATED"
	PaypalSubscriptionStatusInActive PaypalSubscriptionStatus = "INACTIVE"
)

// var PaypalSubscriptionStatus = paypalSubscriptionStatuses{
// 	Active:   "ACTIVE",
// 	Created:  "CREATED",
// 	InActive: "INACTIVE",
// }

type PaypalShowSubscriptionDetailsResponse struct {
	Status PaypalSubscriptionStatus `json:"status"`
}
