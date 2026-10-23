package templates

import (
	"bytes"
)

func ForgotPassword(sessionToken string, verificationToken int32) string {
	var tpl bytes.Buffer
	data := map[string]string{"SessionToken": sessionToken, "VerificationToken": string(verificationToken)} // TODO:complete template data
	emailTemplates.ExecuteTemplate(&tpl, "forgot_password.html", data)
	return tpl.String()
}
