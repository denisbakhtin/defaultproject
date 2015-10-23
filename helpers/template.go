package helpers

import (
	"bytes"
	"html/template"

	"github.com/zenazn/goji/web"
)

func GetTemplate(c web.C) *template.Template {
	return c.Env["Template"].(*template.Template)
}

func Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}
