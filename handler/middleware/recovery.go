package middleware

import (
	"log"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec != nil {
				log.Println("panic is recovered!!!")
			}
		}()
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
