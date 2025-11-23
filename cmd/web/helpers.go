package main

import (
	"net/http"
	"runtime/debug"
)

func (a *app) serverError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Error(err.Error(), "type", "server-error", "method", r.Method, "uri", r.URL.RequestURI(), "trace", string(debug.Stack()))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, r *http.Request, code int, err error) {
	a.logger.Error(err.Error(), "type", "client-error", "method", r.Method, "uri", r.URL.RequestURI())
	http.Error(w, http.StatusText(code), code)
}
