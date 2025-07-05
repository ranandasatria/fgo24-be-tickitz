package dto

type Actor struct {
  ActorName string `json:"actorName" binding:"required"`
}
