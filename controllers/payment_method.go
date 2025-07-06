package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CreatePaymentMethod godoc
// @Summary Create payment method
// @Description Admin only. Add new payment method
// @Tags Payment Method
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreatePaymentMethodRequest true "Payment method name"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/payment-method [post]
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

// GetAllPaymentMethod godoc
// @Summary Get all payment method
// @Description Retrieve all payment method
// @Tags Payment Method
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/payment-method [get]
func GetAllPaymentMethod(c *gin.Context) {
  payment_method, err := models.GetAllPaymentMethod()
  if err != nil {
    c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false, 
			Message: "Failed to fetch payment method", 
			Errors: err.Error()})
    return
  }

  c.JSON(http.StatusOK, utils.Response{
		Success: true, 
		Message: "All payment method", 
		Results: payment_method})
}

// DeletePaymentMethod godoc
// @Summary Delete a payment method
// @Description Admin only. Delete a payment method by ID
// @Tags Payment Method
// @Security BearerAuth
// @Produce json
// @Param id path int true "Payment Method ID"
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/payment-method/{id} [delete]
func DeletePaymentMethod(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	if role, ok := claims["role"].(string); !ok || role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
      Success: false, 
      Message: "Only admin can delete payment method"})
		return
	}

	id := c.Param("id")
	err := models.DeletePaymentMethod(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
      Success: false, 
      Message: "Failed to delete payment method", 
      Errors: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true, 
		Message: "Payment method deleted"})
}