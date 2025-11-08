package main

import (
	"bytes"
	"html/template"
	"net/http"
)

var boardTpl = template.Must(template.ParseGlob("web/templates/*.html"))

type templateData struct {
	EnableTimer   bool
	EnableStandup bool
}

func newTemplateData(c config) templateData {
	return templateData{
		EnableTimer:   c.enableTimer,
		EnableStandup: c.enableStandup,
	}
}

func (a *app) render(w http.ResponseWriter, r *http.Request, status int, data any) {
	// try to render the template, if error return
	buf := new(bytes.Buffer)
	if err := boardTpl.ExecuteTemplate(buf, "base", data); err != nil {
		a.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}
