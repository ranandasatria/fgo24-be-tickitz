package utils

type Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Errors any `json:"error,omitempty"`
	Results any `json:"results,omitempty"`
}