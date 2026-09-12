package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aamirlatif1/imdbapi/internal/data"
	"github.com/aamirlatif1/imdbapi/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	var in input
	err := app.readJSON(w, r, &in)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(in.Title != "", "title", "required")
	v.Check(len(in.Title) <= 500, "title", "must not be larger than 500 bytes")

	v.Check(in.Runtime > 0, "runtime", "must be larger than 0")

	v.Check(in.Year != 0, "year", "must be larger than 0")
	v.Check(in.Year > 1888, "year", "must not be larger than 1888 bytes")

	v.Check(in.Genres != nil, "runtime", "must be provided")
	v.Check(len(in.Genres) > 0, "genres", "must not be empty")
	v.Check(len(in.Genres) <= 5, "genres", "must not be more than 5")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", in)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   120,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
