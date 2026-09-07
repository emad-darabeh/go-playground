package httptransport

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type Response[T any] struct {
	Status int
	Body   T
}

func OK[T any](body T) Response[T] {
	return Response[T]{Status: http.StatusOK, Body: body}
}

func Created[T any](body T) Response[T] {
	return Response[T]{Status: http.StatusCreated, Body: body}
}

func NoContent() Response[struct{}] {
	return Response[struct{}]{Status: http.StatusNoContent}
}

func Handle[T any](fn func(*http.Request) (Response[T], error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := fn(r)
		if err != nil {
			var apiErr Error
			if !errors.As(err, &apiErr) {
				apiErr = InternalServerError()
				slog.Error("unhandled error", "path", r.URL.Path, "err", err)
			}
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(apiErr.Status)
			json.NewEncoder(w).Encode(apiErr)
			return
		}

		if res.Status == http.StatusNoContent {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.Status)
		json.NewEncoder(w).Encode(res.Body)
	}
}
