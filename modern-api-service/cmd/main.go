package main

import (
	"log/slog"
	"net/http"
	"os"

	"modern-api-service/internal/httptransport"
	itemshttp "modern-api-service/internal/items/http"
	"modern-api-service/internal/items"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	itemHandler := itemshttp.NewHandler(items.NewService(items.NewInMemoryStore()))

	srv := httptransport.NewServer(itemHandler)

	addr := ":8080"
	slog.Info("starting server", "addr", addr)
	if err := http.ListenAndServe(addr, srv.Routes()); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
