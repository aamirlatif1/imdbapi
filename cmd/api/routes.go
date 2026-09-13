package main

import (
	"net/http"

	"github.com/aamirlatif1/imdbapi/internal/handlers"
	"github.com/aamirlatif1/imdbapi/internal/store"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) routes(pool *pgxpool.Pool) http.Handler {
	router := mux.NewRouter()

	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedResponse)
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundResponse)

	router.HandleFunc("/v1/health", handlers.HealthHandler).Methods(http.MethodGet)

	handlers.NewMovies(store.NewMovies(pool)).Register(router)

	return handlers.RecoverPanic(handlers.RateLimit(router))
}
