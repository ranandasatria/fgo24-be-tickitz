package controllers

import (
  "be-tickitz/dto"
  "be-tickitz/models"
  "be-tickitz/utils"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v5"
)


// CreateDirector godoc
// @Summary Create director
// @Description Admin only. Add a new director
// @Tags Directors
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.Director true "Director data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/directors [post]
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

// GetAllDirectors godoc
// @Summary Get all directors
// @Description Retrieve all directors
// @Tags Directors
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /directors [get]
func GetAllDirectors(c *gin.Context) {
  directors, err := models.GetAllDirectors()
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{
      Success: false, 
      Message: "Failed to fetch directors", 
      Errors: err.Error(),
    })
    return
  }

  c.JSON(http.StatusOK, utils.Response{
    Success: true, 
    Message: "All directors", 
    Results: directors,
  })
}


// DeleteDirector godoc
// @Summary Delete a director
// @Description Admin only. Delete a director by ID
// @Tags Directors
// @Security BearerAuth
// @Produce json
// @Param id path int true "Director ID"
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/directors/{id} [delete]
func DeleteDirector(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	if role, ok := claims["role"].(string); !ok || role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
      Success: false, 
      Message: "Only admin can delete directors"})
		return
	}

	id := c.Param("id")
	err := models.DeleteDirector(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
      Success: false, 
      Message: "Failed to delete director", 
      Errors: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.Response{Success: true, Message: "Director deleted"})
}