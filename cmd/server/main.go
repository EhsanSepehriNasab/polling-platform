// @title Polling Platform API
// @version 1.0
// @description API for managing polls.
// @host localhost:8080
// @BasePath /
package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/EhsanSepehriNasab/polling-platform/docs"
	"github.com/EhsanSepehriNasab/polling-platform/internal/config"
	"github.com/EhsanSepehriNasab/polling-platform/internal/db"
	"github.com/EhsanSepehriNasab/polling-platform/internal/polls"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.Load()

	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal("Database connection failed:", err)
	}

	db.RunMigrations(cfg.DatabaseURL)

	r := chi.NewRouter()

	// Swagger UI endpoint
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Initialize Poll Repository, Service, and Handler
	pollRepo := polls.NewPollRepository(db.DB)
	pollService := polls.NewPollService(pollRepo)
	pollHandler := polls.NewPollHandler(pollService)

	r.Post("/polls", pollHandler.CreatePoll)

	log.Println("Server starting on port", cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)
	if err != nil {
		log.Fatal(err)
	}
}
