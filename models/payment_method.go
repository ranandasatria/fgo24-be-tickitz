package models

import (
	"be-tickitz/utils"
	"context"

	"github.com/jackc/pgx/v5"
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

func GetAllPaymentMethod() ([]PaymentMethod, error) {
	conn, err := utils.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), `
    SELECT id, payment_name FROM payment_method ORDER BY payment_name ASC
  `)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[PaymentMethod])
}

func DeletePaymentMethod(id string) error {
	conn, err := utils.ConnectDB()
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `DELETE FROM payment_method WHERE id = $1`, id)
	return err
}
