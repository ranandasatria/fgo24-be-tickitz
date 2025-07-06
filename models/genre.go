package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"

	"github.com/jackc/pgx/v5"
)

type Genre struct {
	ID        int    `json:"id"`
	GenreName string `json:"genreName" db:"genre_name"`
}

func CreateGenre(input dto.Genre) (Genre, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return Genre{}, err
	}
	defer conn.Release()

	var genre Genre
	err = conn.QueryRow(context.Background(),`
	INSERT INTO genres (genre_name)
	VALUES ($1)
	RETURNING id, genre_name
	`,
	input.GenreName,
).Scan(
	&genre.ID,
	&genre.GenreName,
)

	return genre, err

}

func AddGenretoMovie(movieID, genreID int) error { //add genre ke movie secara manual melalui hit API
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `
	INSERT INTO movie_genres (id_movie, id_genre)
	VALUES ($1, $2)
	`,
	movieID, genreID,
	)
return err
}

func GetAllGenres() ([]Genre, error){
		conn, err := utils.ConnectDB()
	if err != nil {
		return []Genre{}, err
	}
	defer conn.Release()

	
	rows, err := conn.Query(context.Background(), `
		SELECT id, genre_name
		FROM genres
		ORDER BY genre_name ASC
	`)
	if err != nil {
		return nil, err
	}

	genre,err := pgx.CollectRows(rows, pgx.RowToStructByName[Genre])

	return genre, err
}

func DeleteGenre(id string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `DELETE FROM genres WHERE id = $1`, id)
	return err
}