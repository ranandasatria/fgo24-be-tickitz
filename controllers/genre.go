package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CreateGenre godoc
// @Summary Create genre
// @Description Admin only. Add a new genre
// @Tags Genres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.Genre true "Genre name"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/genres [post]
func CreateGenre(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)

	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Access denied. Only admins can add genres.",
		})
		return
	}

	var genre dto.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
			Errors:  err.Error(),
		})
		return
	}

	created, err := models.CreateGenre(genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create genre",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Genre created successfully.",
		Results: created,
	})
}

func AddGenretoMovie(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)

	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Access denied. Only admins can add a genre to a movie.",
		})
		return
	}

	var req dto.MovieGenres
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	err := models.AddGenretoMovie(req.IDMovie, req.IDGenre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to add genre to movie",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Genre added to movie succesfully",
	})

}

// GetAllGenres godoc
// @Summary Get all genres
// @Description Retrieve all available genres
// @Tags Genres
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /genres [get]
func GetAllGenres(c *gin.Context) {
	genres, err := models.GetAllGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch genres",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "All genres",
		Results: genres,
	})
}

// DeleteGenre godoc
// @Summary Delete a genre
// @Description Admin only. Delete a genre by ID
// @Tags Genres
// @Security BearerAuth
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/genres/{id} [delete]
func DeleteGenre(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	if role, ok := claims["role"].(string); !ok || role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{Success: false, Message: "Only admin can delete genres"})
		return
	}

	id := c.Param("id")
	err := models.DeleteGenre(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{Success: false, Message: "Failed to delete genre", Errors: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.Response{Success: true, Message: "Genre deleted"})
}
