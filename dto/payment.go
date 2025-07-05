package dto

type CreatePaymentMethodRequest struct {
	PaymentName string `json:"paymentName" binding:"required"`
}