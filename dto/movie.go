package dto

import "time"

type Movie struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	ReleaseDate     string `json:"ReleaseDate" db:"release_date"`
	Duration        int    `json:"duration" db:"duration_minutes"`
	Image           string `json:"image"`
	HorizontalImage string `json:"horizontal_image" db:"horizontal_image"`
	GenreIDs        []int  `json:"genreIDs"`
	DirectorIDs     []int  `json:"directorIDs"`
	CastIDs         []int  `json:"castIDs"`
}

type MovieDetail struct {
  ID              int       `json:"id"`
  Title           string    `json:"title"`
  Description     string    `json:"description"`
  ReleaseDate     time.Time `json:"releaseDate"`
  Duration        int       `json:"durationMinutes"`
  Image           string    `json:"image"`
  HorizontalImage string    `json:"horizontalImage"`
  Genres          []string  `json:"genres"`
  Directors       []string  `json:"directors"`
  Casts           []string  `json:"casts"`
}


type MovieResponse struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ReleaseDate     time.Time `json:"releaseDate"`
	Duration        int       `json:"durationMinutes"`
	Image           string    `json:"image"`
	HorizontalImage string    `json:"horizontalImage"`
	GenreIDs        []int     `json:"genreIDs"`
	DirectorIDs     []int     `json:"directorIDs"`
	CastIDs         []int     `json:"castIDs"`
}

type MovieUpcoming struct {
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"ReleaseDate" db:"release_date"`
	Image       string    `json:"image"`
}

type UpdateMovieInput struct {
  Title           *string   `json:"title"`
  Description     *string   `json:"description"`
  ReleaseDate     *string `json:"releaseDate"`
  Duration        *int      `json:"duration"`
  Image           *string   `json:"image"`
  HorizontalImage *string   `json:"horizontalImage"`
  GenreIDs        *[]int    `json:"genreIDs"`
  DirectorIDs     *[]int    `json:"directorIDs"`
  CastIDs         *[]int    `json:"castIDs"`
}

