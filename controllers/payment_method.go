package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CreatePaymentMethod(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	role := claims["role"]

	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Only admin can create payment method",
		})
		return
	}

	var req dto.CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
			Errors:  err.Error(),
		})
		return
	}

	method, err := models.CreatePaymentMethod(req.PaymentName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create payment method",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Payment method created successfully",
		Results: method,
	})
}
