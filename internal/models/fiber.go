package models

const (
	ErrInternalServerError = "oops something went wrong!"
	ErrBadRequest          = "bad request"
	ErrUnauthorized        = "Unauthorized"

	ErrUserNotRegistered  = "user not registered"
	ErrInvalidCredentials = "invalid credentials"
)

type ErrResponse struct {
	Errors  string `json:"errors,omitempty"`
	Message any    `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}
