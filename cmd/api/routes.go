package main

import (
	"net/http"

	"github.com/aamirlatif1/imdbapi/internal/handlers"
	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()

	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedResponse)
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundResponse)

	router.HandleFunc("/v1/health", handlers.HealthHandler).Methods(http.MethodGet)
	movies := handlers.Movies{}
	router.HandleFunc("/v1/movies", movies.Create).Methods(http.MethodPost)
	router.HandleFunc("/v1/movies/{id}", movies.Show).Methods(http.MethodGet)

	return handlers.RecoverPanic(router)
}
