package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"fmt"
	"time"
)

type Movie struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ReleaseDate     time.Time `json:"releaseDate" db:"release_date"`
	Duration        int       `json:"durationMinutes" db:"duration_minutes"`
	Image           string    `json:"image"`
	HorizontalImage string    `json:"horizontalImage" db:"horizontal_image"`
}

func CreateMovie(input dto.Movie) (Movie, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return Movie{}, err
	}
	defer conn.Release()

	parsedDate, err := time.Parse("2006-01-02", input.ReleaseDate)
	if err != nil {
		return Movie{}, fmt.Errorf("invalid release date format: %v", err)
	}

	var movie Movie
	err = conn.QueryRow(context.Background(), `
		INSERT INTO movies (title, description, release_date, duration_minutes, image, horizontal_image)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, release_date, duration_minutes, image, horizontal_image
	`,
		input.Title,
		input.Description,
		parsedDate,
		input.Duration,
		input.Image,
		input.HorizontalImage,
	).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.Duration,
		&movie.Image,
		&movie.HorizontalImage,
	)

	return movie, err
}
