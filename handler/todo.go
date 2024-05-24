package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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
	todoPointers, err := h.svc.ReadTODO(ctx, req.PrevID, int64(req.Size))

	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todoPointers}, nil
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
	err := h.svc.DeleteTODO(ctx, req.IDs)

	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var err error
		len := r.ContentLength
		var readRequest model.ReadTODORequest

		if len != 0 {
			body := make([]byte, len)
			_, err = r.Body.Read(body)

			if err != nil && err.Error() != "EOF" {
				log.Println(err)
			}

			err = json.Unmarshal(body, &readRequest)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			q := r.URL.Query()
			prevID, err := strconv.Atoi(q.Get("prev_id"))
			if err != nil {
				log.Println(err)
			}
			size, err := strconv.Atoi(q.Get("size"))
			if err != nil {
				log.Println(err)
			}

			readRequest.PrevID = int64(prevID)
			if size == 0 {
				readRequest.Size = 5
			} else {
				readRequest.Size = size
			}
		}

		var readResponse *model.ReadTODOResponse

		readResponse, err = h.Read(r.Context(), &readRequest)
		//m, _ := json.Marshal(readResponse)

		if err != nil {
			log.Println(err)
			return
		}
		enc := json.NewEncoder(w)
		enc.Encode(readResponse)

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

	case "DELETE":
		CL := r.ContentLength
		body := make([]byte, CL)

		var DeleteRequest model.DeleteTODORequest

		_, err := r.Body.Read(body)
		if err != nil && err.Error() != "EOF" {
			log.Println(err)
			return
		}

		_ = json.Unmarshal(body, &DeleteRequest)

		if len(DeleteRequest.IDs) == 0 {
			w.WriteHeader(400)
			return
		}

		deleteResponse, err := h.Delete(r.Context(), &DeleteRequest)
		if err != nil {
			log.Println(err.Error())
			if errors.Is(err, &model.ErrNotFound{}) {
				log.Println(err)
				w.WriteHeader(404)
			}
		}

		enc := json.NewEncoder(w)
		enc.Encode(deleteResponse)

	}
}
