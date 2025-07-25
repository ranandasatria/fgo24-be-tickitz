package models

import (
  "be-tickitz/dto"
  "be-tickitz/utils"
  "context"
  "fmt"
  "log"
  "time"
)

func CreateTransaction(userID int, input dto.CreateTransactionRequest) (int, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    return 0, err
  }
  defer conn.Release()

  tx, err := conn.Begin(context.Background())
  if err != nil {
    return 0, err
  }
  defer tx.Rollback(context.Background())

  showDate, err := time.Parse("2006-01-02", input.ShowDate)
  if err != nil {
    return 0, fmt.Errorf("invalid date format: %v", err)
  }

  showTime, err := time.Parse("15:04:05", input.ShowTime)
  if err != nil {
    showTime, err = time.Parse("15:04", input.ShowTime)
    if err != nil {
      return 0, fmt.Errorf("invalid time format: %v", err)
    }
  }

  totalPrice := len(input.Seats) * input.PricePerSeat

  var transactionID int
  err = tx.QueryRow(context.Background(), `
    INSERT INTO transactions (
      id_user, id_movie, show_date, show_time, location, cinema,
      total_price, payment_method
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id
  `, userID, input.MovieID, showDate, showTime, input.Location, input.Cinema, totalPrice, input.PaymentMethod).Scan(&transactionID)

  if err != nil {
    return 0, fmt.Errorf("failed to create transaction: %v", err)
  }

  for _, seat := range input.Seats {
    _, err := tx.Exec(context.Background(), `
      INSERT INTO transaction_details (transaction_id, seat, price)
      VALUES ($1, $2, $3)
    `, transactionID, seat, input.PricePerSeat)

    if err != nil {
      return 0, fmt.Errorf("failed to insert seat %s: %v", seat, err)
    }
  }

  if err := tx.Commit(context.Background()); err != nil {
    return 0, fmt.Errorf("commit failed: %v", err)
  }

  return transactionID, nil
}

func CheckSeatAvailability(input dto.CreateTransactionRequest) ([]string, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    log.Printf("DB connection error: %v", err)
    return nil, err
  }
  defer conn.Release()

  showDate, err := time.Parse("2006-01-02", input.ShowDate)
  if err != nil {
    log.Printf("Error parsing show_date: %v", err)
    return nil, fmt.Errorf("invalid date format: %v", err)
  }

  showTime, err := time.Parse("15:04:05", input.ShowTime)
  if err != nil {
    showTime, err = time.Parse("15:04", input.ShowTime)
    if err != nil {
      log.Printf("Error parsing show_time: %v", err)
      return nil, fmt.Errorf("invalid time format: %v", err)
    }
  }

  query := `
    SELECT td.seat
    FROM transaction_details td
    JOIN transactions t ON td.transaction_id = t.id
    WHERE t.id_movie = $1
      AND t.show_date = $2
      AND t.show_time = $3
      AND LOWER(t.location) = LOWER($4)
      AND LOWER(t.cinema) = LOWER($5)
  `
  params := []interface{}{
    input.MovieID,
    showDate,
    showTime,
    input.Location,
    input.Cinema,
  }

  if len(input.Seats) > 0 {
    query += " AND td.seat = ANY($6)"
    params = append(params, input.Seats)
  }

  rows, err := conn.Query(context.Background(), query, params...)
  if err != nil {
    log.Printf("Query error: %v", err)
    return nil, err
  }
  defer rows.Close()

  var takenSeats []string
  for rows.Next() {
    var seat string
    if err := rows.Scan(&seat); err != nil {
      log.Printf("Scan error: %v", err)
      return nil, err
    }
    takenSeats = append(takenSeats, seat)
  }

  if takenSeats == nil {
    takenSeats = []string{}
  }

  return takenSeats, rows.Err()
}

func GetAllTransactions() ([]dto.TransactionSummary, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    return nil, err
  }
  defer conn.Release()

  rows, err := conn.Query(context.Background(), `
    SELECT
      t.id AS transaction_id,
      m.title AS movie_title,
      t.show_date,
      t.show_time,
      t.location,
      t.cinema,
      t.total_price,
      t.payment_method,
      ARRAY_AGG(td.seat) AS seats
    FROM transactions t
    JOIN movies m ON t.id_movie = m.id
    LEFT JOIN transaction_details td ON td.transaction_id = t.id
    GROUP BY
      t.id, m.title, t.show_date, t.show_time,
      t.location, t.cinema, t.total_price, t.payment_method
    ORDER BY t.created_at DESC
  `)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var transactions []dto.TransactionSummary

  for rows.Next() {
    var t dto.TransactionSummary
    var showDate, showTime time.Time
    var seats []string

    if err := rows.Scan(
      &t.TransactionID,
      &t.MovieTitle,
      &showDate,
      &showTime,
      &t.Location,
      &t.Cinema,
      &t.TotalPrice,
      &t.PaymentMethod,
      &seats,
    ); err != nil {
      return nil, err
    }

    t.ShowDate = showDate.Format("2006-01-02")
    t.ShowTime = showTime.Format("15:04")
    t.Seats = seats

    transactions = append(transactions, t)
  }

  if err := rows.Err(); err != nil {
    return nil, err
  }

  return transactions, nil
}

func GetUserTransactions(userID int) ([]dto.TransactionSummary, error) {
  conn, err := utils.ConnectDB()
  if err != nil {
    return nil, err
  }
  defer conn.Release()

  rows, err := conn.Query(context.Background(), `
    SELECT
      t.id AS transaction_id,
      m.title AS movie_title,
      t.show_date,
      t.show_time,
      t.location,
      t.cinema,
      t.total_price,
      t.payment_method,
      ARRAY_AGG(td.seat) AS seats
    FROM transactions t
    JOIN movies m ON t.id_movie = m.id
    LEFT JOIN transaction_details td ON td.transaction_id = t.id
    WHERE t.id_user = $1
    GROUP BY
      t.id, m.title, t.show_date, t.show_time,
      t.location, t.cinema, t.total_price, t.payment_method
    ORDER BY t.created_at DESC
  `, userID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var transactions []dto.TransactionSummary

  for rows.Next() {
    var t dto.TransactionSummary
    var showDate, showTime time.Time
    var seats []string

    if err := rows.Scan(
      &t.TransactionID,
      &t.MovieTitle,
      &showDate,
      &showTime,
      &t.Location,
      &t.Cinema,
      &t.TotalPrice,
      &t.PaymentMethod,
      &seats,
    ); err != nil {
      return nil, err
    }

    t.ShowDate = showDate.Format("2006-01-02")
    t.ShowTime = showTime.Format("15:04")
    t.Seats = seats

    transactions = append(transactions, t)
  }

  if err := rows.Err(); err != nil {
    return nil, err
  }

  return transactions, nil
}