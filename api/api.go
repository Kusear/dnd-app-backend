package api

import (
	"encoding/json"
	"net/http"
)

type ResponseData[T any] struct {
	// HTTP status code
	Code int
	Data T
}

type ErrorResponse struct {
	Code    int
	Message string
}

func sendErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendResponse[T any](w http.ResponseWriter, code int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	BadRequestError = func(w http.ResponseWriter, message string) {
		sendErrorResponse(w, http.StatusBadRequest, message)
	}
	InternalServerError = func(w http.ResponseWriter) {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
	}
	NotFoundError = func(w http.ResponseWriter) {
		sendErrorResponse(w, http.StatusNotFound, "Not found")
	}
	UnauthorizedError = func(w http.ResponseWriter) {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
	}
	TooManyRequestsError = func(w http.ResponseWriter) {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
	}
)

var (
	SuccessResponse = func(w http.ResponseWriter, data any) {
		sendResponse(w, http.StatusOK, data)
	}
)
