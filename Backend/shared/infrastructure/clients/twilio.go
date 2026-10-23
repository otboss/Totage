package clients

import (
	"waterordersystem/GOBACKEND/utils"

	"github.com/twilio/twilio-go"
	twilioApiClient "github.com/twilio/twilio-go/rest/api/v2010"
)

type twilioApi struct{}

func (t twilioApi) SendSms(recipient string, message string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: utils.Env.TWILIO_ACCOUNT_SID,
		Password: utils.Env.TWILIO_AUTH_TOKEN,
	})

	params := &twilioApiClient.CreateMessageParams{}
	params.SetTo(recipient)
	params.SetFrom(utils.Env.TWILIO_PHONE_NUMBER)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	return nil
}

var TwilioApi = twilioApi{}
