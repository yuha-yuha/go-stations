package handler

import (
	"context"
	"encoding/json"
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
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
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
		var createTodoRequest model.CreateTODORequest

		_, err := r.Body.Read(body)
		log.Println(string(body))
		if err != nil {
			log.Println(err, "57")
		}

		err = json.Unmarshal(body, &createTodoRequest)
		log.Println("log:::", createTodoRequest.Description)
		if err != nil {
			log.Println(err, "64")
			return
		}

		if createTodoRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte{})
			return
		}

		var dbResult *model.CreateTODOResponse
		dbResult, _ = h.Create(r.Context(), &createTodoRequest)
		enc := json.NewEncoder(w)
		enc.Encode(dbResult)
	}
}
