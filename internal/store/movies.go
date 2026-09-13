package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aamirlatif1/imdbapi/internal/validator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Movie struct {
	ID        int64     `json:"id"`               // Unique integer ID for the movie
	CreatedAt time.Time `json:"-"`                // Timestamp for when the movie is added to our database
	Title     string    `json:"title"`            // Movie title
	Year      int32     `json:"year,omitzero"`    // Movie release year
	Runtime   int32     `json:"runtime,omitzero"` // Movie runtime (in minutes)
	Genres    []string  `json:"genres,omitzero"`  // Slice of genres for the movie (romance, comedy, etc.)
	Version   int32     `json:"version"`          // The version number starts at 1 and will be incremented each
}

const movieColumns = "id, title, year, runtime, genres, created_at, version"

type Movies struct {
	db *pgxpool.Pool
}

func NewMovies(db *pgxpool.Pool) *Movies {
	return &Movies{
		db: db,
	}
}

func (m Movies) Add(ctx context.Context, movie *Movie) (*Movie, error) {
	const query = `INSERT INTO movies (title, year, runtime, genres)
					VALUES ($1, $2, $3, $4)
					RETURNING ` + movieColumns

	row := m.db.QueryRow(ctx, query, movie.Title, movie.Year, movie.Runtime, movie.Genres)

	out, err := scanMovie(row)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("movie already exists")
		}
		return nil, fmt.Errorf("adding movie failed: %w", err)
	}

	return &out, nil
}

func (m Movies) Get(ctx context.Context, id int64) (*Movie, error) {
	query := `SELECT ` + movieColumns + ` FROM movies WHERE id = $1`

	movie, err := scanMovie(m.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("movie not found")
	}
	if err != nil {
		return nil, fmt.Errorf("getting movie failed: %w", err)
	}
	return &movie, nil
}

func (m Movies) List(ctx context.Context, limit, offset int32) ([]Movie, error) {
	return nil, nil
}

func (m Movies) Update(ctx context.Context, movie *Movie) (*Movie, error) {
	query := `update movies set title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
				where id = $5
				returning ` + movieColumns
	row := m.db.QueryRow(ctx, query, movie.Title, movie.Year, movie.Runtime, movie.Genres, movie.ID)
	out, err := scanMovie(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.New("movie not found")
		}
		return nil, fmt.Errorf("updating movie failed: %w", err)
	}
	return &out, nil
}

func (m Movies) Delete(ctx context.Context, id int64) error {
	query := `delete from movies where id = $1`
	result, err := m.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting movie failed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("movie not found")
	}
	return nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "required")
	v.Check(len(movie.Title) <= 500, "title", "must not be larger than 500 bytes")

	v.Check(movie.Runtime > 0, "runtime", "must be larger than 0")

	v.Check(movie.Year != 0, "year", "must be larger than 0")
	v.Check(movie.Year > 1888, "year", "must not be larger than 1888 bytes")

	v.Check(movie.Genres != nil, "runtime", "must be provided")
	v.Check(len(movie.Genres) > 0, "genres", "must not be empty")
	v.Check(len(movie.Genres) <= 5, "genres", "must not be more than 5")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicates")
}

type scanner interface {
	Scan(dest ...any) error
}

func scanMovie(row scanner) (Movie, error) {
	var m Movie
	err := row.Scan(&m.ID, &m.Title, &m.Year, &m.Runtime, &m.Genres, &m.CreatedAt, &m.Version)
	return m, err
}
