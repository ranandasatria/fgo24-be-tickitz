package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CreateMovie godoc
// @Summary Create new movie
// @Description Admin only. Add a new movie with metadata and relations
// @Tags Movies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.Movie true "Movie data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/movies [post]
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
		DirectorIDs:     movie.DirectorIDs,
		CastIDs:         movie.CastIDs,
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Movie created successfully",
		Results: response,
	})
}

// GetAllMovies godoc
// @Summary Get all movies
// @Description View all movies in database
// @Tags Movies
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /movies [get]
func GetAllMovies(c *gin.Context) {
	ctx := context.Background()

	err := utils.RedisClient().Ping(ctx).Err()
	noredis := false
	if err != nil {
		log.Println("Redis unavailable:", err.Error())
		noredis = true
	}

	search := c.DefaultQuery("search", "")
	cacheKey := "/movies?search=" + search

	if !noredis {
		val, err := utils.RedisClient().Get(ctx, cacheKey).Result()
		if err == nil {
			var cached []dto.MovieList
			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				c.JSON(http.StatusOK, utils.Response{
					Success: true,
					Message: "All movies (from Redis)",
					Results: cached,
				})
				return
			}
		}
	}

	var rawMovies []models.Movie
	if search == "" {
		rawMovies, err = models.GetAllMovies()
	} else {
		rawMovies, err = models.SearchMovies(search)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch movies",
			Errors:  err.Error(),
		})
		return
	}

	var movies []dto.MovieList
	for _, m := range rawMovies {
		movies = append(movies, dto.MovieList{
			ID:              m.ID,
			Title:           m.Title,
			Description:     m.Description,
			ReleaseDate:     m.ReleaseDate,
			Duration:        m.Duration,
			Image:           m.Image,
			HorizontalImage: m.HorizontalImage,
		})
	}

	if !noredis {
		if encoded, err := json.Marshal(movies); err == nil {
			utils.RedisClient().Set(ctx, cacheKey, encoded, 0)
		}
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "All movies",
		Results: movies,
	})
}

// GetMovieByID godoc
// @Summary Get movie by ID
// @Description Retrieve movie details by its ID
// @Tags Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /movies/{id} [get]
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

// GetNowShowing godoc
// @Summary Get now showing movies
// @Description Retrieve list of currently showing movies
// @Tags Movies
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /movies/now-showing [get]
func GetNowShowing(c *gin.Context) {
	movies, err := models.GetNowShowing()
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
		Message: "Now showing",
		Results: movies,
	})
}

// GetUpcoming godoc
// @Summary Get upcoming movies
// @Description Retrieve list of upcoming movies
// @Tags Movies
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /movies/upcoming [get]
func GetUpcoming(c *gin.Context) {
	movies, err := models.GetUpcoming()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch movies",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Upcoming movies",
		Results: movies,
	})
}

// DeleteMovie godoc
// @Summary Delete a movie
// @Description Admin only. Delete a movie by ID
// @Tags Movies
// @Security BearerAuth
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/movies/{id} [delete]
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

// UpdateMovie godoc
// @Summary Update a movie
// @Description Admin only. Update movie details and relations
// @Tags Movies
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param request body dto.UpdateMovieInput true "Update movie data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/movies/{id} [patch]
func UpdateMovie(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	role, _ := claims["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Only admin can update movie",
		})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var input dto.UpdateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
			Errors:  err.Error(),
		})
		return
	}

	if err := models.UpdateMovie(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to update movie",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Movie updated successfully",
	})
}
