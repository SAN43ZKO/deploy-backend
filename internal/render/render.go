package render

import (
	"encoding/json"
	"net/http"
)

const (
	ExpiredToken = iota + 1
)

type Err struct {
	Message string `json:"message" validate:"required"`
	Code    *int   `json:"code,omitempty" validate:"required"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func Error(w http.ResponseWriter, status int, message string, code ...int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var errCode *int
	if len(code) > 0 {
		errCode = &code[0]
	}

	data := Err{
		Message: message,
		Code:    errCode,
	}

	json.NewEncoder(w).Encode(data)
}
