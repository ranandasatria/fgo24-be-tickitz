package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CreateMovie(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)

	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Access denied. Only admins can add movies.",
		})
		return
	}

	var movie dto.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
			Errors:  err.Error(),
		})
		return
	}

	created, err := models.CreateMovie(movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create movie",
			Errors:  err.Error(),
		})
		return
	}

	

	response := dto.MovieResponse{
  ID:              created.ID,
  Title:           created.Title,
  Description:     created.Description,
  ReleaseDate:     created.ReleaseDate,
  Duration:        created.Duration,
  Image:           created.Image,
  HorizontalImage: created.HorizontalImage,
  GenreIDs:        movie.GenreIDs,
}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Movie created successfully",
		Results: response,
	})
}

func GetAllMovies(c *gin.Context) {
	movies, err := models.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch movies",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "All movies",
		Results: movies,
	})
}

func GetMovieByID(c *gin.Context) {
	id := c.Param("id")

	movie, err := models.GetMovieByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.Response{
			Success: false,
			Message: "Movie not found",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Movie details",
		Results: movie,
	})
}

func GetNowShowing(c *gin.Context){
	movies, err := models.GetNowShowing()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch movies",
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Now showing",
		Results: movies,
	})
}

func GetUpcoming(c *gin.Context){
	movies, err := models.GetUpcoming()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch movies",
			Errors: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Upcoming movies",
		Results: movies,
	})
}

func DeleteMovie(c *gin.Context) {
  claims := c.MustGet("user").(jwt.MapClaims)
  if role, ok := claims["role"].(string); !ok || role != "admin" {
    c.JSON(http.StatusForbidden, utils.Response{Success: false, Message: "Only admin can delete movies"})
    return
  }

  id := c.Param("id")
  err := models.DeleteMovie(id)
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{Success: false, Message: "Failed to delete movie", Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{Success: true, Message: "Movie deleted"})
}