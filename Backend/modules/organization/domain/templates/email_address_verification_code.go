package templates

import (
	"bytes"
)

func EmailAddressVerificationCode(sessionToken string, verificationToken int32) string {
	var tpl bytes.Buffer
	data := map[string]string{"SessionToken": sessionToken, "VerificationToken": string(verificationToken)} // TODO:complete template data
	emailTemplates.ExecuteTemplate(&tpl, "email_address_verification_code.html", data)
	return tpl.String()
}
