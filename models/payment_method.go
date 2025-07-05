package models

import (
	"be-tickitz/utils"
	"context"
)

type PaymentMethod struct {
	ID          int    `json:"id"`
	PaymentName string `json:"paymentName"`
}

func CreatePaymentMethod(paymentName string) (PaymentMethod, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return PaymentMethod{}, err
	}
	defer conn.Release()

	var method PaymentMethod
	err = conn.QueryRow(context.Background(), `
		INSERT INTO payment_method (payment_name)
		VALUES ($1)
		RETURNING id, payment_name
	`, paymentName).Scan(
		&method.ID,
		&method.PaymentName,
	)

	return method, err
}
