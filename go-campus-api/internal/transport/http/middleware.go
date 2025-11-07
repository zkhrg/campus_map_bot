package http

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func withMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func methodValidator(allowedMethod string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != allowedMethod {
				respondError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func jsonValidator() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")
			if contentType != "application/json" {
				respondError(w, http.StatusUnsupportedMediaType, "content-type must be application/json")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
