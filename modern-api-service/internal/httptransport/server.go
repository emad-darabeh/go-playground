package httptransport

import (
	"log/slog"
	"net/http"
)

type Handler interface {
	Register(mux *http.ServeMux)
}

type Server struct {
	handlers []Handler
}

func NewServer(handlers ...Handler) *Server {
	return &Server{handlers: handlers}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	for _, h := range s.handlers {
		h.Register(mux)
	}
	return withLogging(mux)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
