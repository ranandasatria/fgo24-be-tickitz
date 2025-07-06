package controllers

import (
  "be-tickitz/dto"
  "be-tickitz/models"
  "be-tickitz/utils"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)

// CreateActors godoc
// @Summary Create actor
// @Description Admin only. Add a new actor
// @Tags Actors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.Actor true "Actor data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/actors [post]
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

// GetAllActors godoc
// @Summary Get all actors
// @Description Retrieve all actors
// @Tags Actors
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/actors [get]
func GetAllActors(c *gin.Context) {
  actors, err := models.GetAllActors()
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{Success: false, Message: "Failed to fetch actors", Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{Success: true, Message: "All actors", Results: actors})
}


// DeleteActor godoc
// @Summary Delete a actor
// @Description Admin only. Delete a actor by ID
// @Tags Actors
// @Security BearerAuth
// @Produce json
// @Param id path int true "Actor ID"
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/actors/{id} [delete]
func DeleteActor(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	if role, ok := claims["role"].(string); !ok || role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
      Success: false, 
      Message: "Only admin can delete actors"})
		return
	}

	id := c.Param("id")
	err := models.DeleteActor(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
      Success: false, 
      Message: "Failed to delete actor", 
      Errors: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.Response{Success: true, Message: "Actor deleted"})
}