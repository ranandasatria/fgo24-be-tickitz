package models

import (
  "be-tickitz/dto"
  "be-tickitz/utils"
  "context"

  "github.com/jackc/pgx/v5"
)

type Actor struct {
  ID         int    `json:"id"`
  ActorName  string `json:"actorName" db:"actor_name"`
}

func CreateActor(input dto.Actor) (Actor, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    return Actor{}, err
  }
  defer conn.Release()

  var actor Actor
  err = conn.QueryRow(context.Background(), `
    INSERT INTO actors (actor_name)
    VALUES ($1)
    RETURNING id, actor_name
  `, input.ActorName).Scan(&actor.ID, &actor.ActorName)

  return actor, err
}

func GetAllActors() ([]Actor, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    return nil, err
  }
  defer conn.Release()

  rows, err := conn.Query(context.Background(), `
    SELECT id, actor_name FROM actors ORDER BY actor_name ASC
  `)
  if err != nil {
    return nil, err
  }

  return pgx.CollectRows(rows, pgx.RowToStructByName[Actor])
}

func DeleteActor(id string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `DELETE FROM actors WHERE id = $1`, id)
	return err
}