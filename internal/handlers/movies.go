package handlers

import (
	"context"
	"net/http"
	"time"

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

func (m *Movies) Create(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, http.StatusCreated, saved, nil)
}

func (m *Movies) Show(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie := store.Movie{
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

func (m *Movies) Register(router *mux.Router) {
	router.HandleFunc("/v1/movies", m.Create).Methods(http.MethodPost)
	router.HandleFunc("/v1/movies/{id}", m.Show).Methods(http.MethodGet)
}
