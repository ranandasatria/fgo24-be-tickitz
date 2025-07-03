package models

import (
	"be-tickitz/dto"
	"be-tickitz/utils"
	"context"
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