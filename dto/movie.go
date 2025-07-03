package dto

import "time"

type Movie struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	ReleaseDate     string `json:"ReleaseDate" db:"release_date"`
	Duration        int    `json:"duration" db:"duration_minutes"`
	Image           string `json:"image"`
	HorizontalImage string `json:"horizontal_image" db:"horizontal_image"`
}

type MovieUpcoming struct {
	Title       string  `json:"title"`
	ReleaseDate time.Time `json:"ReleaseDate" db:"release_date"`
	Image       string  `json:"image"`
}
