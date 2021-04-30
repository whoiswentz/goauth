package app

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAppError(w http.ResponseWriter, e error, statusCode int) []byte {
	appErrorBytes, _ := json.Marshal(AppError{
		Code:    statusCode,
		Message: e.Error(),
	})

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return appErrorBytes
}
