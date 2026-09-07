package httptransport

import (
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	Type     string `json:"type,omitempty"`
	Status   int    `json:"status"`
	Title    string `json:"title"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func (e Error) Error() string {
	parts := []string{fmt.Sprintf("status=%d title=%s", e.Status, e.Title)}
	if e.Type != "" {
		parts = append(parts, fmt.Sprintf("type=%s", e.Type))
	}
	if e.Detail != "" {
		parts = append(parts, fmt.Sprintf("detail=%s", e.Detail))
	}
	if e.Instance != "" {
		parts = append(parts, fmt.Sprintf("instance=%s", e.Instance))
	}
	return strings.Join(parts, " ")
}

func NotFoundError(detail string) Error {
	return Error{
		Type:   "/problems/not-found",
		Status: http.StatusNotFound,
		Title:  http.StatusText(http.StatusNotFound),
		Detail: detail,
	}
}

func BadRequestError(detail string) Error {
	return Error{
		Type:   "/problems/bad-request",
		Status: http.StatusBadRequest,
		Title:  http.StatusText(http.StatusBadRequest),
		Detail: detail,
	}
}

func InternalServerError() Error {
	return Error{
		Type:   "/problems/internal-server",
		Status: http.StatusInternalServerError,
		Title:  http.StatusText(http.StatusInternalServerError),
	}
}
