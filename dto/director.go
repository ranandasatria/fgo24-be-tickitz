package dto

type Director struct {
  DirectorName string `json:"directorName" binding:"required"`
}