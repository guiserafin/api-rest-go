package main

import (
	"api/api"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {

	if err := run(); err != nil {
		slog.Error("Failed to run", "error", err)
		os.Exit(1)
	}

	slog.Info("System offline")

}

func run() error {

	handler := api.HttpHandler()

	s := http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
