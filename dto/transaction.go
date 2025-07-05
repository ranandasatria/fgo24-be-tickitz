package dto

type CreateTransactionRequest struct {
  MovieID       int      `json:"movie_id"`
  ShowDate      string   `json:"show_date"` 
  ShowTime      string   `json:"show_time"` 
  Location      string   `json:"location"`
  Cinema        string   `json:"cinema"`
  Seats         []string `json:"seats"`         
  PricePerSeat  int      `json:"price_per_seat"`
  PaymentMethod int      `json:"payment_method"` 
}
