package middleware

import "net/http"

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer recover()
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
