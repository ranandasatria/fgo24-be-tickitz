package controllers

import (
  "be-tickitz/dto"
  "be-tickitz/models"
  "be-tickitz/utils"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

func CreateDirector(c *gin.Context) {
  claims := c.MustGet("user").(jwt.MapClaims)
  if role, ok := claims["role"].(string); !ok || role != "admin" {
    c.JSON(http.StatusForbidden, utils.Response{Success: false, Message: "Only admin can add directors"})
    return
  }

  var input dto.Director
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, utils.Response{Success: false, Message: "Invalid input", Errors: err.Error()})
    return
  }

  director, err := models.CreateDirector(input)
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{Success: false, Message: "Failed to create director", Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{Success: true, Message: "Director created", Results: director})
}

func GetAllDirectors(c *gin.Context) {
  directors, err := models.GetAllDirectors()
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{Success: false, Message: "Failed to fetch directors", Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{Success: true, Message: "All directors", Results: directors})
}
