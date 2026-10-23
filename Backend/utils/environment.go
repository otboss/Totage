package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	Env.IS_PRODUCTION = strings.ToLower(os.Getenv("IS_PRODUCTION")) == "true"
	Env.IS_ORGANIZATION_SUBSCRIPTION_ENABLED = strings.ToLower(os.Getenv("IS_ORGANIZATION_SUBSCRIPTION_ENABLED")) == "true"
	Env.BACKEND_URL = os.Getenv("BACKEND_URL")

	Env.DATABASE_HOST = os.Getenv("DATABASE_HOST")
	Env.DATABASE_NAME = os.Getenv("DATABASE_NAME")
	Env.DATABASE_USER = os.Getenv("DATABASE_USER")
	Env.DATABASE_PASSWORD = os.Getenv("DATABASE_PASSWORD")
	Env.DATABASE_SSL_ENABLED = strings.ToLower(os.Getenv("DATABASE_SSL_ENABLED")) == "true"
	Env.CAPTCHA_DISABLED = strings.ToLower(os.Getenv("CAPTCHA_DISABLED")) == "true"
	Env.GOOGLE_CLOUD_API_KEY = os.Getenv("GOOGLE_CLOUD_API_KEY")

	Env.BTC_ELECTRUM_SEED_PHRASE = os.Getenv("BTC_ELECTRUM_SEED_PHRASE")
	Env.BTC_LIGHTNING_SEED = os.Getenv("BTC_LIGHTNING_SEED")

	Env.KAFKA_GROUP_ID = os.Getenv("KAFKA_GROUP_ID")
	Env.KAFKA_BOOTSTRAP_SERVERS = os.Getenv("KAFKA_BOOTSTRAP_SERVERS")

	Env.TWILIO_ACCOUNT_SID = os.Getenv("TWILIO_ACCOUNT_SID")
	Env.TWILIO_PHONE_NUMBER = os.Getenv("TWILIO_PHONE_NUMBER")
	Env.TWILIO_AUTH_TOKEN = os.Getenv("TWILIO_AUTH_TOKEN")

	Env.EMAIL_DISABLED = strings.ToLower(os.Getenv("EMAIL_DISABLED")) == "true"
	Env.SMTP_HOST = os.Getenv("SMTP_HOST")
	Env.SMTP_PORT = os.Getenv("SMTP_PORT")
	Env.SMTP_USER = os.Getenv("SMTP_USER")
	Env.SMTP_NAME = os.Getenv("SMTP_NAME")
	Env.SMTP_PASS = os.Getenv("SMTP_PASS")

	Env.OAUTH_GOOGLE_CLIENT_ID = os.Getenv("OAUTH_CLIENT_ID")
	Env.OAUTH_GOOGLE_CLIENT_SECRET = os.Getenv("OAUTH_CLIENT_SECRET")

	Env.FRONTEND_ORIGIN = os.Getenv("FRONTEND_ORIGIN")
	Env.BACKEND_ORIGIN = os.Getenv("BACKEND_ORIGIN")

	Env.ENC_KEY = os.Getenv("ENC_KEY")

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	Env.PORT = uint16(port)

	dbPort, err := strconv.ParseUint(os.Getenv("DATABASE_PORT"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	Env.DATABASE_PORT = uint16(dbPort)

	distributionCenterOperationalRadius, err := strconv.ParseFloat(os.Getenv("DISTRIBUTION_CENTER_OPERATIONAL_RADIUS_KM"), 32)
	if err != nil {
		log.Fatalln(err)
	}
	Env.DISTRIBUTION_CENTER_OPERATIONAL_RADIUS_KM = float32(distributionCenterOperationalRadius)

	distributionCenterCommissionPercent, err := strconv.ParseFloat(os.Getenv("DISTRIBUTION_CENTER_COMMISSION_PERCENT"), 32)
	if err != nil {
		log.Fatalln(err)
	}
	Env.DISTRIBUTION_CENTER_COMMISSION_PERCENT = float32(distributionCenterCommissionPercent)

	Env.DISTRIBUTION_CENTER_DEFAULT_ID = os.Getenv("DISTRIBUTION_CENTER_DEFAULT_ID")
}

type environment struct {
	IS_PRODUCTION                        bool
	IS_ORGANIZATION_SUBSCRIPTION_ENABLED bool
	BACKEND_URL                          string
	PORT                                 uint16

	BTC_ELECTRUM_SEED_PHRASE string
	BTC_LIGHTNING_SEED       string

	DATABASE_NAME        string
	DATABASE_HOST        string
	DATABASE_PORT        uint16
	DATABASE_USER        string
	DATABASE_PASSWORD    string
	DATABASE_SSL_ENABLED bool

	PAYPAL_REST_SDK_MODE          string
	PAYPAL_REST_SDK_CLIENT_ID     string
	PAYPAL_REST_SDK_CLIENT_SECRET string
	PAYPAL_API_BASE_URL           string

	CURRENCYAPI_KEY          string
	EXCHANGE_RATE_USD_TO_JMD float32

	TWILIO_ACCOUNT_SID  string
	TWILIO_PHONE_NUMBER string
	TWILIO_AUTH_TOKEN   string

	CAPTCHA_DISABLED     bool
	GOOGLE_CLOUD_API_KEY string

	DISTRIBUTION_CENTER_OPERATIONAL_RADIUS_KM float32
	DISTRIBUTION_CENTER_COMMISSION_PERCENT    float32
	DISTRIBUTION_CENTER_DEFAULT_ID            string

	KAFKA_GROUP_ID          string
	KAFKA_BOOTSTRAP_SERVERS string

	EMAIL_DISABLED bool
	SMTP_HOST      string
	SMTP_NAME      string
	SMTP_PORT      string
	SMTP_USER      string
	SMTP_PASS      string

	OAUTH_GOOGLE_CLIENT_ID     string
	OAUTH_GOOGLE_CLIENT_SECRET string

	FRONTEND_ORIGIN string
	BACKEND_ORIGIN  string

	ENC_KEY string
}

var Env *environment = &environment{}
