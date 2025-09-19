package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SuccessResponse defines standard api success response
type SuccessResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
	Success bool   `json:"success"`
}

// ErrorResponse defines standard api error response
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

// RespondWithError returns consistent error messages
func RespondWithError(w http.ResponseWriter, statusCode int, message string, err string) {
	resp := ErrorResponse{
		Message: message,
		Error:   err,
		Success: false,
	}

	fmt.Println("Error response: ", message, " - ", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// RespondWithJSON sends a JSON response with the given status code, message, and data
func RespondWithJSON(w http.ResponseWriter, statusCode int, message string, data any) error {
	resp := SuccessResponse[any]{
		Message: message,
		Data:    data,
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(resp)
}
