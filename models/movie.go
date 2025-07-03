package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return Movie{}, err
	}

	defer tx.Rollback(context.Background())

	parsedDate, err := time.Parse("2006-01-02", input.ReleaseDate)
	if err != nil {
		return Movie{}, fmt.Errorf("invalid release date format: %v", err)
	}

	var movie Movie
	err = tx.QueryRow(context.Background(), `
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
	if err != nil{
		return Movie{}, err
	}

	for _, genreID := range input.GenreIDs {
		_, err := tx.Exec(context.Background(), `
		INSERT INTO movie_genres (id_movie, id_genre)
		VALUES ($1, $2)
		`, movie.ID, genreID)

		if err != nil {
			fmt.Printf("Insert genre failed: %v\n", err)
			return Movie{}, fmt.Errorf("failed to insert genre: %v", err)
		}
	}
	if err := tx.Commit(context.Background()); err != nil {
  return Movie{}, fmt.Errorf("failed to commit transaction: %v", err)
	}
	return movie, nil
}

func GetAllMovies() ([]Movie, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
		SELECT id, title, description, release_date, duration_minutes, image, horizontal_image
		FROM movies
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Movie])
	return movies, err
}

func GetMovieByID(id string) (Movie, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return Movie{}, err
	}
	defer conn.Release()

	var movie Movie
	err = conn.QueryRow(context.Background(), `
		SELECT id, title, description, release_date, duration_minutes, image, horizontal_image
		FROM movies
		WHERE id = $1
	`, id).Scan(
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


func GetNowShowing() ([]Movie, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	
	rows, err := conn.Query(context.Background(), `
		SELECT id, title, description, release_date, duration_minutes, image, horizontal_image
		FROM movies
		WHERE release_date < CURRENT_DATE
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Movie])
	return movies, err
}

func GetUpcoming() ([]dto.MovieUpcoming, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
	SELECT title, release_date, image
	FROM movies
	WHERE release_date > CURRENT_DATE
	ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.MovieUpcoming])
	return movies, err
}