package app

type AppError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
