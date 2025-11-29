package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/ekaputra07/go-retro/web/ui"
)

var boardTpl = template.Must(template.ParseFS(ui.UiFS, "dist/*.html"))

type templateData struct {
	AppName    string
	AppVersion string
	AppTagline string
}

type templateAndJSONData struct {
	templateData
	JSONData template.JS
}

func newTemplateData(_ config) (*templateAndJSONData, error) {
	data := templateData{
		AppName:    appName,
		AppVersion: appVersion,
		AppTagline: appTagline,
	}
	// create JSON string version of the data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &templateAndJSONData{
		templateData: data,
		JSONData:     template.JS(jsonData),
	}, nil
}

func (a *app) render(w http.ResponseWriter, r *http.Request, status int, data any) {
	// try to render the template, if error return
	buf := new(bytes.Buffer)
	if err := boardTpl.ExecuteTemplate(buf, "index.html", data); err != nil {
		a.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}
