package templates

import (
	"bytes"
	db_order "waterordersystem/GOBACKEND/modules/order/infrastructure"
	db_user "waterordersystem/GOBACKEND/modules/user/infrastructure"
)

func OrderPlaced(customer db_user.GetUserRow, order db_order.GetOrderByIdRow) string {
	var tpl bytes.Buffer
	data := map[string]string{} // TODO:complete template data
	emailTemplates.ExecuteTemplate(&tpl, "order_placed.html", data)
	return tpl.String()
}
