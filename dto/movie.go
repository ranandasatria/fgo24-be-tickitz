package dto

type Movie struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"ReleaseDate" db:"release_date"`
	Duration int `json:"duration" db:"duration_minutes"`
	Image string `json:"image"`
	HorizontalImage string `json:"horizontal_image" db:"horizontal_image"`
}
