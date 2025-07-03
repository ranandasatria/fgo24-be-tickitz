package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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
