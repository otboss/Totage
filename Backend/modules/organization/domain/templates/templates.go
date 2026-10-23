package templates

import (
	"embed"
	"html/template"
)

//go:embed email_templates/*.html
var emailTemplatesFS embed.FS

var emailTemplates = template.Must(template.ParseFS(
	emailTemplatesFS,
	"email_templates/*.html",
))
