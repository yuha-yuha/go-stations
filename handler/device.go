package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

type DeviceHandler struct{}

func NewDeviceHandler() http.Handler {
	return DeviceHandler{}
}

func (h DeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ResponseBody struct {
		OS string `json:"os"`
	}

	ResponseBody.OS, _ = r.Context().Value(model.DeviceKey("OS")).(string)

	enc := json.NewEncoder(w)

	enc.Encode(ResponseBody)

}
