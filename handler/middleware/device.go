package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

func Device(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		type DeviceKey string

		ua := useragent.Parse(r.UserAgent())
		uaos := ua.OS

		os := DeviceKey("OS")

		ctx := context.WithValue(r.Context(), os, uaos)

		h.ServeHTTP(w, r.WithContext(ctx))

	}

	return http.HandlerFunc(fn)
}
