package main

import (
	"fmt"
	"net/http"
)

// commonHeaders set headers that always set for every response
func (a *app) commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'self';")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Xss-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

// staticCache set cache-control header for all static files
func (a *app) staticCache(next http.Handler, maxAgeSecond int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAgeSecond))
		next.ServeHTTP(w, r)
	})
}
