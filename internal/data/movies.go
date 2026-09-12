package data

import (
	"time"

	"github.com/aamirlatif1/imdbapi/internal/validator"
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
