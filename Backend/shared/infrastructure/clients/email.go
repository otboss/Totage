package clients

import (
	"net/smtp"
	"waterordersystem/GOBACKEND/utils"
)

type emailApi struct{}

type sendEmailParams struct {
	From          string `validate:"required,regexp=^.+\\s<.+@.+>$"`
	To            []string
	Subject       string
	HtmlEmailBody string
}

func (n emailApi) SendEmail(sender string, recipients []string, subject string, body string) error {
	msg := []byte(subject + "\n" + body)
	auth := smtp.PlainAuth("", sender, utils.Env.SMTP_PASS, utils.Env.SMTP_HOST)

	return smtp.SendMail(
		utils.Env.SMTP_HOST+":"+utils.Env.SMTP_PORT,
		auth,
		sender,
		recipients,
		msg,
	)
}

var EmailApi = emailApi{}
