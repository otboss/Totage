package clients

type currencyApi struct {
}

var exchangeRateUsdToJmd float32

// func init() {
// 	currencyApi := currencyApi{}

// 	go func() {
// 		for {
// 			response, err := currencyApi.exchangeRateUsdToJmd()
// 			if err != nil {
// 				if exchangeRateUsdToJmd == 0 {
// 					log.Fatalln(err.Error())
// 				} else {
// 					log.Println(err.Error())
// 				}
// 			}
// 			exchangeRateUsdToJmd = float32(response.Data.JMD.Value)
// 			time.Sleep(12 * time.Hour)
// 		}
// 	}()
// }

// func (api currencyApi) exchangeRateUsdToJmd() (*currencyApiResponse, error) {
// 	// Construct the URL with API key
// 	url := fmt.Sprintf("https://api.currencyapi.com/v3/latest?apikey=%s&currencies=JMD", utils.Env.CURRENCYAPI_KEY)

// 	// Perform the HTTP GET request
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// Read and parse the response
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var response currencyApiResponse
// 	if err := json.Unmarshal(body, &response); err != nil {
// 		return nil, err
// 	}

// 	return &response, nil
// }

func (api currencyApi) GetExchangeRateUsd(currencyCodeAlpha3 string) (float64, error) {
	return 1, nil
}

func (api currencyApi) GetExchangeRateUsdToJmd() float32 {
	return exchangeRateUsdToJmd
}

var CurrencyApi = currencyApi{}

type currencyApiResponse struct {
	Data struct {
		JMD struct {
			Value float64 `json:"value"`
		} `json:"JMD"`
	} `json:"data"`
}
