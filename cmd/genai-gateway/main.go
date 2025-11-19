package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := getEnv("HTTP_PORT", "8094")
	adress := ":" + port

	logger := log.New(os.Stdout, "[genai-gateway] ", log.LstdFlags|log.Lshortfile)

	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)

	// Middlewares propios
	////
	////

	service := &http.Server{
		Addr:         adress,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Print("HTTP server listening on %s\n", adress)
		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %v\n", err)
		}
	}()
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
