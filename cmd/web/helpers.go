package main

import "net/http"

func (a *app) serverError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, r *http.Request, code int, err error) {
	a.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
	http.Error(w, http.StatusText(code), code)
}
