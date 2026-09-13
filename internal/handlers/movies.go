package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	err = writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (m *Movies) Update(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Title   *string  `json:"title"`
		Year    *int32   `json:"year"`
		Runtime *int32   `json:"runtime"`
		Genres  []string `json:"genres"`
	}
	id, err := readIDParam(r)

	if err != nil {
		NotFoundResponse(w, r)
		return
	}

	movie, err := m.store.Get(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			NotFoundResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(int(movie.Version)) != r.Header.Get("X-Expected-Version") {
			editConflictResponse(w, r)
			return
		}
	}

	var in input
	err = readJSON(w, r, &in)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if in.Title != nil {
		movie.Title = *in.Title
	}
	if in.Year != nil {
		movie.Year = *in.Year
	}
	if in.Runtime != nil {
		movie.Runtime = *in.Runtime
	}
	if in.Genres != nil {
		movie.Genres = in.Genres
	}

	v := validator.New()
	if store.ValidateMovie(v, movie); !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}
	saved, err := m.store.Update(r.Context(), movie)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEditConflict):
			editConflictResponse(w, r)
		default:
			serverErrorResponse(w, r, err)
		}
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", saved.ID))
	err = writeJSON(w, http.StatusOK, saved, headers)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (m *Movies) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := readIDParam(r)
	if err != nil {
		NotFoundResponse(w, r)
	}
	err = m.store.Delete(r.Context(), id)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	err = writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

func (m *Movies) Register(router *mux.Router) {
	router.HandleFunc("/v1/movies", m.Create).Methods(http.MethodPost)
	router.HandleFunc("/v1/movies/{id}", m.Show).Methods(http.MethodGet)
	router.HandleFunc("/v1/movies/{id}", m.Update).Methods(http.MethodPatch)
	router.HandleFunc("/v1/movies/{id}", m.Delete).Methods(http.MethodDelete)
}
