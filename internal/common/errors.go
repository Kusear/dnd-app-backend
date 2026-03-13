package common

type ErrorType string

const (
	BadRequestError      ErrorType = "bad_request"
	InternalServerError  ErrorType = "internal_server_error"
	NotFoundError        ErrorType = "not_found"
	UnauthorizedError    ErrorType = "unauthorized"
	TooManyRequestsError ErrorType = "too_many_requests"
)

func (e ErrorType) GetErrorMessage(message string) map[string]any {
	return map[string]any{
		"type":    e,
		"message": message,
	}
}
