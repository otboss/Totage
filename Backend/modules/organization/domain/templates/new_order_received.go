package templates

import (
	"bytes"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
)

func NewOrderReceived(customer db_user.GetUserRow, order db_order.GetOrderByIdRow) string {
	var tpl bytes.Buffer
	data := map[string]string{} // TODO:complete template data
	emailTemplates.ExecuteTemplate(&tpl, "new_order_received.html", data)
	return tpl.String()
}
