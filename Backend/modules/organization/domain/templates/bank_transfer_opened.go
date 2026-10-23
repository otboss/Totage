package templates

import (
	"bytes"
)

func BankTransferOpened() string {
	var tpl bytes.Buffer
	data := map[string]string{} // TODO:complete template data
	emailTemplates.ExecuteTemplate(&tpl, "forgot_password.html", data)
	return tpl.String()
}
