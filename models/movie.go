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
	GenreIDs        []int     `json:"genre_ids"`
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

	for _, genreID := range input.GenreIDs {
		_, err := tx.Exec(context.Background(), `
      INSERT INTO movie_genres (id_movie, id_genre)
      VALUES ($1, $2)
    `, movie.ID, genreID)
		if err != nil {
			return Movie{}, fmt.Errorf("failed to insert genre: %v", err)
		}
	}

	for _, directorID := range input.DirectorIDs {
		_, err := tx.Exec(context.Background(), `
      INSERT INTO movie_directors (id_movie, id_director)
      VALUES ($1, $2)
    `, movie.ID, directorID)
		if err != nil {
			return Movie{}, fmt.Errorf("failed to insert director: %v", err)
		}
	}

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


	utils.DeleteKeysByPrefix(context.Background(), "/movies")

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
    ORDER BY id ASC
  `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.Duration, &m.Image, &m.HorizontalImage)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func SearchMovies(search string) ([]Movie, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `
    SELECT id, title, description, release_date, duration_minutes, image, horizontal_image
    FROM movies
    WHERE title ILIKE '%' || $1 || '%'
    ORDER BY id ASC
  `
	rows, err := conn.Query(context.Background(), query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.Duration, &m.Image, &m.HorizontalImage)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func GetMovieByID(id string) (dto.MovieDetail, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return dto.MovieDetail{}, err
	}
	defer conn.Release()

	var movie dto.MovieDetail
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
	if err != nil {
		return dto.MovieDetail{}, err
	}

	rows, err := conn.Query(context.Background(), `
    SELECT g.genre_name
    FROM genres g
    JOIN movie_genres mg ON g.id = mg.id_genre
    WHERE mg.id_movie = $1
  `, id)
	if err == nil {
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err == nil {
				movie.Genres = append(movie.Genres, name)
			}
		}
	}

	rows, err = conn.Query(context.Background(), `
    SELECT d.director_name
    FROM directors d
    JOIN movie_directors md ON d.id = md.id_director
    WHERE md.id_movie = $1
  `, id)
	if err == nil {
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err == nil {
				movie.Directors = append(movie.Directors, name)
			}
		}
	}

	rows, err = conn.Query(context.Background(), `
    SELECT a.actor_name
    FROM actors a
    JOIN movie_casts mc ON a.id = mc.id_actor
    WHERE mc.id_movie = $1
  `, id)
	if err == nil {
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err == nil {
				movie.Casts = append(movie.Casts, name)
			}
		}
	}

	return movie, nil
}

func GetNowShowing(search string, genres []int, sort string, page int, limit int) ([]Movie, int, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, 0, err
	}
	defer conn.Release()

	// Query utama untuk movies
	query := `
    SELECT id, title, description, release_date, duration_minutes, image, horizontal_image
    FROM movies
    WHERE release_date <= NOW()
  `
	countQuery := `
    SELECT COUNT(id)
    FROM movies
    WHERE release_date <= NOW()
  `
	params := []interface{}{}

	if search != "" {
		query += ` AND title ILIKE '%' || $1 || '%'`
		countQuery += ` AND title ILIKE '%' || $1 || '%'`
		params = append(params, search)
	}

	if len(genres) > 0 {
		query += fmt.Sprintf(` AND id IN (
      SELECT id_movie FROM movie_genres WHERE id_genre = ANY($%d)
    )`, len(params)+1)
		countQuery += fmt.Sprintf(` AND id IN (
      SELECT id_movie FROM movie_genres WHERE id_genre = ANY($%d)
    )`, len(params)+1)
		params = append(params, genres)
	}

	if sort == "latest" {
		query += ` ORDER BY release_date DESC`
	} else if sort == "name-asc" {
		query += ` ORDER BY title ASC`
	} else if sort == "name-desc" {
		query += ` ORDER BY title DESC`
	} else {
		query += ` ORDER BY id ASC`
	}

	if page > 0 && limit > 0 {
		query += fmt.Sprintf(` OFFSET $%d LIMIT $%d`, len(params)+1, len(params)+2)
		params = append(params, (page-1)*limit, limit)
	}

	// Hitung total rows
	var totalRows int
	err = conn.QueryRow(context.Background(), countQuery, params[:len(params)-2]...).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}

	// Ambil movies
	rows, err := conn.Query(context.Background(), query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.Duration, &m.Image, &m.HorizontalImage)
		if err != nil {
			return nil, 0, err
		}
		movies = append(movies, m)
	}

	// Ambil genre_ids untuk setiap movie
	for i, m := range movies {
		var genreIDs []int
		rows, err := conn.Query(context.Background(), `
      SELECT id_genre
      FROM movie_genres
      WHERE id_movie = $1
    `, m.ID)
		if err == nil {
			for rows.Next() {
				var id int
				if err := rows.Scan(&id); err == nil {
					genreIDs = append(genreIDs, id)
				}
			}
			rows.Close()
		}
		movies[i].GenreIDs = genreIDs
	}

	return movies, totalRows, nil
}

func GetUpcoming() ([]Movie, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
    SELECT m.id, m.title, m.description, m.release_date, m.duration_minutes, m.image, m.horizontal_image,
           COALESCE(ARRAY_AGG(mg.id_genre) FILTER (WHERE mg.id_genre IS NOT NULL), '{}') as genre_ids
    FROM movies m
    LEFT JOIN movie_genres mg ON m.id = mg.id_movie
    WHERE m.release_date > NOW()
    GROUP BY m.id
  `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Movie])
	return movies, err
}

func DeleteMovie(id string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `DELETE FROM movies WHERE id = $1`, id)
	if err == nil {
		utils.DeleteKeysByPrefix(context.Background(), "/movies")
	}
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
	err = tx.QueryRow(context.Background(), `
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
		parsedDate, err := time.Parse("2006-01-02", *input.ReleaseDate)
		if err != nil {
			return fmt.Errorf("invalid release date format: %v", err)
		}
		releaseDate = parsedDate
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

	if input.DirectorIDs != nil {
		_, err := tx.Exec(context.Background(), `DELETE FROM movie_directors WHERE id_movie = $1`, id)
		if err != nil {
			return err
		}
		for _, directorID := range *input.DirectorIDs {
			_, err := tx.Exec(context.Background(), `
        INSERT INTO movie_directors (id_movie, id_director)
        VALUES ($1, $2)
      `, id, directorID)
			if err != nil {
				return err
			}
		}
	}

	if input.CastIDs != nil {
		_, err := tx.Exec(context.Background(), `DELETE FROM movie_casts WHERE id_movie = $1`, id)
		if err != nil {
			return err
		}
		for _, actorID := range *input.CastIDs {
			_, err := tx.Exec(context.Background(), `
        INSERT INTO movie_casts (id_movie, id_actor)
        VALUES ($1, $2)
      `, id, actorID)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit(context.Background())
	if err == nil {
		utils.DeleteKeysByPrefix(context.Background(), "/movies")
	}
	return err
}
