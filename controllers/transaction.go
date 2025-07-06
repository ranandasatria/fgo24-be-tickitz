package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"

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
// @Security BearerAuth
// @Param request body dto.CreateTransactionRequest true "Transaction request body"
// @Success 200 {object} utils.Response
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

  err = models.CreateTransaction(userID, input)
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
  })
}

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
      Message: "Only admin can access"})
    return
  }

  results, err := models.GetAllTransactions()
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{
      Success: false, 
      Message: "Failed to fetch", 
      Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{
    Success: true,
    Message: "All transactions",
    Results: results,
  })
}

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
      Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{
    Success: true,
    Message: "Your transactions",
    Results: transactions,
  })
}