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
  if err != nil {
    return Movie{}, err
  }

  // insert genre IDs
  for _, genreID := range input.GenreIDs {
    _, err := tx.Exec(context.Background(), `
      INSERT INTO movie_genres (id_movie, id_genre)
      VALUES ($1, $2)
    `, movie.ID, genreID)

    if err != nil {
      return Movie{}, fmt.Errorf("failed to insert genre: %v", err)
    }
  }

  // insert director IDs
  for _, directorID := range input.DirectorIDs {
    _, err := tx.Exec(context.Background(), `
      INSERT INTO movie_directors (id_movie, id_director)
      VALUES ($1, $2)
    `, movie.ID, directorID)
    if err != nil {
      return Movie{}, fmt.Errorf("failed to insert director: %v", err)
    }
  }

  // insert cast IDs
  for _, actorID := range input.CastIDs {
    _, err := tx.Exec(context.Background(), `
      INSERT INTO movie_casts (id_movie, id_actor)
      VALUES ($1, $2)
    `, movie.ID, actorID)
    if err != nil {
      return Movie{}, fmt.Errorf("failed to insert cast: %v", err)
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

func DeleteMovie(id string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `DELETE FROM movies WHERE id = $1`, id)
	return err
}

func UpdateMovie(id int, input dto.UpdateMovieInput) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var old Movie
	err = conn.QueryRow(context.Background(), `
		SELECT id, title, description, release_date, duration_minutes, image, horizontal_image
		FROM movies WHERE id = $1
	`, id).Scan(
		&old.ID,
		&old.Title,
		&old.Description,
		&old.ReleaseDate,
		&old.Duration,
		&old.Image,
		&old.HorizontalImage,
	)
	if err != nil {
		return fmt.Errorf("movie not found: %v", err)
	}

	title := old.Title
	if input.Title != nil {
		title = *input.Title
	}

	description := old.Description
	if input.Description != nil {
		description = *input.Description
	}

	releaseDate := old.ReleaseDate
	if input.ReleaseDate != nil {
		releaseDate = *input.ReleaseDate
	}

	duration := old.Duration
	if input.Duration != nil {
		duration = *input.Duration
	}

	image := old.Image
	if input.Image != nil {
		image = *input.Image
	}

	horizontalImage := old.HorizontalImage
	if input.HorizontalImage != nil {
		horizontalImage = *input.HorizontalImage
	}

	_, err = tx.Exec(context.Background(), `
		UPDATE movies SET 
			title = $1, 
			description = $2, 
			release_date = $3,
			duration_minutes = $4, 
			image = $5, 
			horizontal_image = $6,
			updated_at = NOW()
		WHERE id = $7
	`, title, description, releaseDate, duration, image, horizontalImage, id)

	if err != nil {
		return err
	}

	if input.GenreIDs != nil {
		_, err := tx.Exec(context.Background(), `DELETE FROM movie_genres WHERE id_movie = $1`, id)
		if err != nil {
			return err
		}

		for _, genreID := range *input.GenreIDs {
			_, err := tx.Exec(context.Background(), `
				INSERT INTO movie_genres (id_movie, id_genre)
				VALUES ($1, $2)
			`, id, genreID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(context.Background())
}
