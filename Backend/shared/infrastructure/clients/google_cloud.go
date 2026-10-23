package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"waterordersystem/GOBACKEND/utils"
)

type googleCloudApi struct {
}

type distanceMatrixApiResponse struct {
	DestinationAddresses, OriginAddresses []string `json:"destination_addresses,omitempty",json:"origin_addresses,omitempty"`
	Rows                                  []struct {
		Elements []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

// geocodeApiResponse is the response structure for the Geocode API
// type GeocodeApiResponse struct {
// 	Results []struct {
// 		AddressComponents []struct {
// 			LongName, ShortName string   `json:"long_name, short_name"`
// 			Types               []string `json:"types"`
// 		} `json:"address_components"`
// 		FormattedAddress string `json:"formatted_address"`
// 		Geometry         struct {
// 			Bounds struct {
// 				Northeast struct {
// 					Lat, Lng float64 `json:"lat, lng"`
// 				} `json:"northeast"`
// 				Southwest struct {
// 					Lat, Lng float64 `json:"lat, lng"`
// 				} `json:"southwest"`
// 			} `json:"bounds"`
// 			Location struct {
// 				Lat, Lng float64 `json:"lat, lng"`
// 			} `json:"location"`
// 			LocationType string `json:"location_type"`
// 			Viewport     struct {
// 				Northeast2 struct {
// 					Lat, Lng float64 `json:"lat, lng"`
// 				} `json:"northeast"`
// 				Southwest2 struct {
// 					Lat, Lng float64 `json:"lat, lng"`
// 				} `json:"southwest"`
// 			} `json:"viewport"`
// 		} `json:"geometry"`
// 		PlaceID string   `json:"place_id"`
// 		Types   []string `json:"types"`
// 	} `json:"results"`
// 	Status string `json:"status"`
// }

type GeocodeApiResponse struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Bounds struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"bounds"`
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
			Viewport     struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

// recaptchaEventRequest defines the request body for reCAPTCHA validation
type RecaptchaEventRequest struct {
	Event RecaptchaEvent `json:"event"`
}

type RecaptchaEvent struct {
	Token          string `json:"token"`
	ExpectedAction string `json:"expectedAction"`
	SiteKey        string `json:"siteKey"`
}

// RecaptchaEventResponseDto defines the response for reCAPTCHA validation
type recaptchaEventResponseDto struct {
	Name  string `json:"name"`
	Event struct {
		Token                    string   `json:"token"`
		SiteKey                  string   `json:"siteKey"`
		UserAgent                string   `json:"userAgent"`
		UserIpAddress            string   `json:"userIpAddress"`
		ExpectedAction           string   `json:"expectedAction"`
		HashedAccountId          string   `json:"hashedAccountId"`
		Express                  bool     `json:"express"`
		RequestedUri             string   `json:"requestedUri"`
		WafTokenAssessment       bool     `json:"wafTokenAssessment"`
		Ja3                      string   `json:"ja3"`
		Headers                  []string `json:"headers"`
		FirewallPolicyEvaluation bool     `json:"firewallPolicyEvaluation"`
		FraudPrevention          string   `json:"fraudPrevention"`
	} `json:"event"`
	RiskAnalysis struct {
		Score                  float64       `json:"score"`
		Reasons                []interface{} `json:"reasons"`
		ExtendedVerdictReasons []interface{} `json:"extendedVerdictReasons"`
	} `json:"riskAnalysis"`
	TokenProperties struct {
		Valid              bool   `json:"valid"`
		InvalidReason      string `json:"invalidReason"`
		Hostname           string `json:"hostname"`
		AndroidPackageName string `json:"androidPackageName"`
		IosBundleId        string `json:"iosBundleId"`
		Action             string `json:"action"`
		CreateTime         string `json:"createTime"`
	} `json:"tokenProperties"`
}

func (g googleCloudApi) ReverseGeocoding(latitude float64, longitude float64) (reverseGeocodingResponse, error) {
	// TODO:
	// https://maps.googleapis.com/maps/api/geocode/json?latlng=17.971,-76.792&key=YOUR_API_KEY
	return reverseGeocodingResponse{}, nil
}

// distanceBetweenLocations gets the distance between two locations using the Distance Matrix API
func (g googleCloudApi) DistanceBetweenLocations(origin, destination Location) (*distanceMatrixApiResponse, error) {
	origin.Address = url.QueryEscape(origin.Address)
	destination.Address = url.QueryEscape(destination.Address)

	originCoordinates := origin.Address
	destinationCoordinates := destination.Address

	if origin.Coordinates != nil {
		originCoordinates = fmt.Sprintf("%.7f,%.7f", origin.Coordinates.Longitude, origin.Coordinates.Latitude)
	}

	if destination.Coordinates != nil {
		destinationCoordinates = fmt.Sprintf("%.7f,%.7f", destination.Coordinates.Longitude, destination.Coordinates.Latitude)
	}

	apiURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&units=imperial&key=%s", originCoordinates, destinationCoordinates, utils.Env.GOOGLE_CLOUD_API_KEY)

	// Send the GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var response distanceMatrixApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

type Location struct {
	Address     string
	Coordinates *Coordinates
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// geocodeAddress geocodes an address using the Geocoding API
func (g googleCloudApi) GeocodeAddress(address string) (*GeocodeApiResponse, error) {
	address = url.QueryEscape(address)

	apiURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", address, utils.Env.GOOGLE_CLOUD_API_KEY)

	// Send the GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response GeocodeApiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (g googleCloudApi) OAuthTokenInfo(accessToken string) (*tokenInfo, error) {
	apiURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?access_token=%s", accessToken)

	// Send the GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response tokenInfo
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type tokenInfo struct {
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Scope         string `json:"scope"`
	Exp           string `json:"exp"`
	ExpiresIn     string `json:"expires_in"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AccessType    string `json:"access_type"`
}

func (g googleCloudApi) OAuthUserInfo(accessToken string) (*googleUserInfo, error) {
	apiURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accessToken)

	// Send the GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response googleUserInfo
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// recaptchaValidation validates the reCAPTCHA event using the reCAPTCHA API
func (g googleCloudApi) RecaptchaValidation(recaptchaEvent RecaptchaEventRequest) (*recaptchaEventResponseDto, error) {
	// Convert recaptchaEvent to JSON
	eventBody, err := json.Marshal(recaptchaEvent)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("https://recaptchaenterprise.googleapis.com/v1/projects/waterordersystem-425218/assessments?key=%s", utils.Env.GOOGLE_CLOUD_API_KEY)

	// Send the POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(eventBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var response recaptchaEventResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (g googleCloudApi) OAuth(google_oauth_authorization_code string) (string, error) {
	// TODO: Complete OAUTH function
	return "", nil
}

var GoogleCloudApi = googleCloudApi{}

type reverseGeocodingResponse struct {
	Status  string `json:"status"`
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"` // <--- This field holds the Country Code (e.g., "US")
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			LocationType string `json:"location_type"`
		} `json:"geometry"`
		PlaceID string   `json:"place_id"`
		Types   []string `json:"types"`
	} `json:"results"`
}
