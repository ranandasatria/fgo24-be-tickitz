package models

import (
  "be-tickitz/dto"
  "be-tickitz/utils"
  "context"

  "github.com/jackc/pgx/v5"
)

type Director struct {
  ID           int    `json:"id"`
  DirectorName string `json:"directorName" db:"director_name"`
}

func CreateDirector(input dto.Director) (Director, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    return Director{}, err
  }
  defer conn.Release()

  var director Director
  err = conn.QueryRow(context.Background(), `
    INSERT INTO directors (director_name)
    VALUES ($1)
    RETURNING id, director_name
  `, input.DirectorName).Scan(&director.ID, &director.DirectorName)

  return director, err
}

func GetAllDirectors() ([]Director, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    return nil, err
  }
  defer conn.Release()

  rows, err := conn.Query(context.Background(), `
    SELECT id, director_name FROM directors ORDER BY director_name ASC
  `)
  if err != nil {
    return nil, err
  }

  return pgx.CollectRows(rows, pgx.RowToStructByName[Director])
}

func DeleteDirector(id string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `DELETE FROM directors WHERE id = $1`, id)
	return err
}