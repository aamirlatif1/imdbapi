package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	router.HandleFunc("/v1/health", app.heathHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/movies", app.createMovieHandler).Methods(http.MethodPost)
	router.HandleFunc("/v1/movies/{id}", app.showMovieHandler).Methods(http.MethodGet)

	return app.recoverPanic(router)
}
