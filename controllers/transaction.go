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
