package itemshttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"modern-api-service/internal/httptransport"
	"modern-api-service/internal/items"
)

type itemDTO struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func toDTO(item items.Item) itemDTO {
	return itemDTO{ID: item.ID, Name: item.Name}
}

type Handler struct {
	service *items.Service
}

func NewHandler(service *items.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /items", httptransport.Handle(h.list))
	mux.HandleFunc("POST /items", httptransport.Handle(h.create))
	mux.HandleFunc("GET /items/{id}", httptransport.Handle(h.get))
	mux.HandleFunc("DELETE /items/{id}", httptransport.Handle(h.delete))
}

func (h *Handler) list(r *http.Request) (httptransport.Response[[]itemDTO], error) {
	all := h.service.List()
	dtos := make([]itemDTO, len(all))
	for i, item := range all {
		dtos[i] = toDTO(item)
	}
	return httptransport.OK(dtos), nil
}

func (h *Handler) create(r *http.Request) (httptransport.Response[itemDTO], error) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return httptransport.Response[itemDTO]{}, httptransport.BadRequestError("invalid body")
	}

	item, err := h.service.Create(body.Name)
	if err != nil {
		return httptransport.Response[itemDTO]{}, mapError(err)
	}

	slog.Info("item created", "id", item.ID, "name", item.Name)
	return httptransport.Created(toDTO(item)), nil
}

func (h *Handler) get(r *http.Request) (httptransport.Response[itemDTO], error) {
	id, err := parseID(r)
	if err != nil {
		return httptransport.Response[itemDTO]{}, err
	}

	item, err := h.service.Get(id)
	if err != nil {
		return httptransport.Response[itemDTO]{}, mapError(err)
	}
	return httptransport.OK(toDTO(item)), nil
}

func (h *Handler) delete(r *http.Request) (httptransport.Response[struct{}], error) {
	id, err := parseID(r)
	if err != nil {
		return httptransport.Response[struct{}]{}, err
	}

	if err := h.service.Delete(id); err != nil {
		return httptransport.Response[struct{}]{}, mapError(err)
	}

	slog.Info("item deleted", "id", id)
	return httptransport.NoContent(), nil
}

// mapError translates domain errors to HTTP transport errors.
func mapError(err error) httptransport.Error {
	switch {
	case errors.Is(err, items.ErrNotFound):
		return httptransport.NotFoundError(err.Error())
	case errors.Is(err, items.ErrNameRequired):
		return httptransport.BadRequestError(err.Error())
	default:
		return httptransport.InternalServerError()
	}
}
