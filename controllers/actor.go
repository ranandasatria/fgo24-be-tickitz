package controllers

import (
  "be-tickitz/dto"
  "be-tickitz/models"
  "be-tickitz/utils"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

func CreateActor(c *gin.Context) {
  claims := c.MustGet("user").(jwt.MapClaims)
  if role, ok := claims["role"].(string); !ok || role != "admin" {
    c.JSON(http.StatusForbidden, utils.Response{Success: false, Message: "Only admin can add actors"})
    return
  }

  var input dto.Actor
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, utils.Response{Success: false, Message: "Invalid input", Errors: err.Error()})
    return
  }

  actor, err := models.CreateActor(input)
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{Success: false, Message: "Failed to create actor", Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{Success: true, Message: "Actor created", Results: actor})
}

func GetAllActors(c *gin.Context) {
  actors, err := models.GetAllActors()
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{Success: false, Message: "Failed to fetch actors", Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{Success: true, Message: "All actors", Results: actors})
}
