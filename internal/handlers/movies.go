package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aamirlatif1/imdbapi/internal/data"
	"github.com/aamirlatif1/imdbapi/internal/validator"
)

type Movies struct {
}

func (app *Movies) Create(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	var in input
	err := readJSON(w, r, &in)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   in.Title,
		Year:    in.Year,
		Runtime: in.Runtime,
		Genres:  in.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", in)
}

func (app *Movies) Show(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)

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
	err = writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
