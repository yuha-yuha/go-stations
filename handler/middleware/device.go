package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mileusna/useragent"
)

func Device(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ua := useragent.Parse(r.UserAgent())
		log.Println(r.UserAgent())
		uaos := ua.OS

		os := model.DeviceKey("OS")

		ctx := context.WithValue(r.Context(), os, uaos)

		h.ServeHTTP(w, r.WithContext(ctx))

	}

	return http.HandlerFunc(fn)
}
