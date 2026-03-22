package mail

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/*.tmpl
var mailTemplates embed.FS

type codeMailData struct {
	Code string
}

// VerifyEmailHTML — как templates/verify-email.hbs в Nest.
func VerifyEmailHTML(code string) (string, error) {
	return renderMailTmpl("templates/verify-email.tmpl", codeMailData{Code: code})
}

// ForgotPasswordHTML — как templates/forgot-password.hbs в Nest.
func ForgotPasswordHTML(code string) (string, error) {
	return renderMailTmpl("templates/forgot-password.tmpl", codeMailData{Code: code})
}

func renderMailTmpl(path string, data any) (string, error) {
	raw, err := mailTemplates.ReadFile(path)
	if err != nil {
		return "", err
	}
	t, err := template.New(path).Parse(string(raw))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
