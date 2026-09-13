package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aamirlatif1/imdbapi/internal/store"
	"github.com/aamirlatif1/imdbapi/internal/validator"
	"github.com/gorilla/mux"
)

type MovieStore interface {
	Add(ctx context.Context, movie *store.Movie) (*store.Movie, error)
	Get(ctx context.Context, id int64) (*store.Movie, error)
	List(ctx context.Context, limit, offset int32) ([]store.Movie, error)
	Update(ctx context.Context, movie *store.Movie) (*store.Movie, error)
	Delete(ctx context.Context, id int64) error
}
type Movies struct {
	store MovieStore
}

func NewMovies(s MovieStore) *Movies {
	return &Movies{
		store: s,
	}
}

type input struct {
	Title   string   `json:"title"`
	Year    int32    `json:"year"`
	Runtime int32    `json:"runtime"`
	Genres  []string `json:"genres"`
}

func (m *Movies) Create(w http.ResponseWriter, r *http.Request) {
	var in input
	err := readJSON(w, r, &in)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	movie := &store.Movie{
		Title:   in.Title,
		Year:    in.Year,
		Runtime: in.Runtime,
		Genres:  in.Genres,
	}

	v := validator.New()

	if store.ValidateMovie(v, movie); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}
	saved, err := m.store.Add(r.Context(), movie)
	if err != nil {
		badRequestResponse(w, r, err)
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", saved.ID))
	err = writeJSON(w, http.StatusCreated, saved, headers)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (m *Movies) Show(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)

	if err != nil {
		NotFoundResponse(w, r)
		return
	}

	movie, err := m.store.Get(r.Context(), id)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, movie, nil)
}

func (m *Movies) Update(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)

	if err != nil {
		NotFoundResponse(w, r)
		return
	}

	movie, err := m.store.Get(r.Context(), id)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	var in input
	err = readJSON(w, r, &in)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	movie.Runtime = in.Runtime
	movie.Title = in.Title
	movie.Year = in.Year
	movie.Genres = in.Genres

	v := validator.New()
	if store.ValidateMovie(v, movie); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
	}
	saved, err := m.store.Update(r.Context(), movie)
	if err != nil {
		badRequestResponse(w, r, err)
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", saved.ID))
	err = writeJSON(w, http.StatusOK, saved, headers)
}

func (m *Movies) Register(router *mux.Router) {
	router.HandleFunc("/v1/movies", m.Create).Methods(http.MethodPost)
	router.HandleFunc("/v1/movies/{id}", m.Show).Methods(http.MethodGet)
	router.HandleFunc("/v1/movies/{id}", m.Update).Methods(http.MethodPut)
}
