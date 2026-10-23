package templates

import (
	"bytes"
)

func OrganizationInvitation(accountId string) string {
	var tpl bytes.Buffer
	data := map[string]string{} // TODO:complete template data
	emailTemplates.ExecuteTemplate(&tpl, "order_placed.html", data)
	return tpl.String()
}
