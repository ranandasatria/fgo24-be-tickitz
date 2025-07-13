package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Book seats and create a new transaction
// @Tags Transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Transaction request body"
// @Success 200 {object} utils.Response{results=int}
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /transactions [post]
func CreateTransaction(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	userID := int(claims["userId"].(float64))

	var input dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
			Errors:  err.Error(),
		})
		return
	}

	timeRegex := regexp.MustCompile(`([0-1][0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?`)
	if !strings.Contains(input.ShowTime, ":") || !timeRegex.MatchString(input.ShowTime) {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid show_time format",
		})
		return
	}

	conflictSeats, err := models.CheckSeatAvailability(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to check seat availability",
			Errors:  err.Error(),
		})
		return
	}

	if len(conflictSeats) > 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Seats are already taken",
			Results: conflictSeats,
		})
		return
	}

	transactionID, err := models.CreateTransaction(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create transaction",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Transaction created successfully",
		Results: transactionID,
	})
}

// CheckSeatAvailability godoc
// @Summary Check seat availability
// @Description Check if seats are available for a movie show
// @Tags Transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param movie_id query int true "Movie ID"
// @Param show_date query string true "Show date (YYYY-MM-DD)"
// @Param show_time query string true "Show time (HH:MM:SS)"
// @Param location query string true "Location"
// @Param cinema query string true "Cinema"
// @Param seats query string false "Seats (comma-separated)"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /check-seats [get]
func CheckSeatAvailability(c *gin.Context) {
	var input dto.CreateTransactionRequest
	movieIDStr := c.Query("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil || movieID <= 0 {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid movie_id",
		})
		return
	}
	input.MovieID = movieID
	input.ShowDate = c.Query("show_date")
	input.ShowTime = c.Query("show_time")
	input.Location = c.Query("location")
	input.Cinema = c.Query("cinema")
	seats := c.Query("seats")
	if seats != "" {
		input.Seats = strings.Split(seats, ",")
	} else {
		input.Seats = []string{}
	}

	timeRegex := regexp.MustCompile(`([0-1][0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?`)
	if !strings.Contains(input.ShowTime, ":") || !timeRegex.MatchString(input.ShowTime) {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid show_time format",
		})
		return
	}

	takenSeats, err := models.CheckSeatAvailability(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to check seat availability",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Seat availability checked",
		Results: takenSeats,
	})
}

// GetAllTransactions godoc
// @Summary Get all transactions (admin only)
// @Tags Transactions
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/transactions [get]
func GetAllTransactions(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	role := claims["role"].(string)

	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Only admin can access",
		})
		return
	}

	results, err := models.GetAllTransactions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "All transactions",
		Results: results,
	})
}

// GetMyTransactions godoc
// @Summary Get transactions per user
// @Tags Transactions
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /transactions [get]
func GetMyTransactions(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	userID := int(claims["userId"].(float64))

	transactions, err := models.GetUserTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Your transactions",
		Results: transactions,
	})
}
