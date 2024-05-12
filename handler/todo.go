package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)

	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	result, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *result}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Println("POST")
		len := r.ContentLength
		body := make([]byte, len)
		var createRequest model.CreateTODORequest

		_, err := r.Body.Read(body)
		log.Println(string(body))
		if err != nil && err.Error() != "EOF" {
			log.Println(err)
		}

		err = json.Unmarshal(body, &createRequest)
		log.Println("log:::", createRequest.Description)
		if err != nil {
			log.Println(err)
			return
		}

		if createRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		var dbResult *model.CreateTODOResponse
		dbResult, _ = h.Create(r.Context(), &createRequest)
		enc := json.NewEncoder(w)
		enc.Encode(dbResult)

	case "PUT":
		len := r.ContentLength
		body := make([]byte, len)

		var updateRequest model.UpdateTODORequest

		_, err := r.Body.Read(body)
		if err != nil && err.Error() != "EOF" {
			log.Println(err)
			return
		}

		err = json.Unmarshal(body, &updateRequest)
		if err != nil {
			log.Println(err)
			return
		}

		if updateRequest.ID == 0 || updateRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
		}

		var updateResponse *model.UpdateTODOResponse

		updateResponse, err = h.Update(r.Context(), &updateRequest)

		if err != nil {
			if errors.Is(err, model.ErrNotFound{}) {
				log.Println(err)
				w.WriteHeader(http.StatusNotFound)
			}

			log.Println(err)
			return
		}

		enc := json.NewEncoder(w)
		enc.Encode(updateResponse)
	}
}
