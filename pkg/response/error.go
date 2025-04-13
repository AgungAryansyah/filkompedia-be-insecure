package response

import (
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Err  error `json:"error"`
	Code int   `json:"code"`
}

func (e *ErrorResponse) Error() string {
	return e.Err.Error()
}

func NewErrorResponse(code int, message string) ErrorResponse {
	return ErrorResponse{
		Code: code,
		Err:  errors.New(message),
	}
}

var (
	InternalServerError = NewErrorResponse(http.StatusInternalServerError, "Internal server error")

	UserNotFound     = NewErrorResponse(http.StatusNotFound, "User not found")
	CommentNotFound  = NewErrorResponse(http.StatusNotFound, "Comment not found")
	BookNotFound     = NewErrorResponse(http.StatusNotFound, "Book not found")
	CartNotFound     = NewErrorResponse(http.StatusNotFound, "Cart not found")
	CheckoutNotFound = NewErrorResponse(http.StatusNotFound, "Checkout not found")
	PaymentNotFound  = NewErrorResponse(http.StatusNotFound, "Payment not found")

	UserUnverified   = NewErrorResponse(http.StatusForbidden, "User is not verified")
	DuplicateAccount = NewErrorResponse(http.StatusConflict, "User already exists")

	BadRequest = NewErrorResponse(http.StatusBadRequest, "Bad request")

	InvalidToken       = NewErrorResponse(http.StatusUnauthorized, "Token invalid")
	InvalidOTP         = NewErrorResponse(http.StatusUnauthorized, "OTP invalid")
	ExpiredToken       = NewErrorResponse(http.StatusUnauthorized, "Expired token")
	InvalidCredentials = NewErrorResponse(http.StatusUnauthorized, "Invalid credentials")

	Unauthorized     = NewErrorResponse(http.StatusUnauthorized, "Unauthorized access")
	RoleUnauthorized = NewErrorResponse(http.StatusForbidden, "Insufficient role")
	Forbidden        = NewErrorResponse(http.StatusForbidden, "Forbidden access")
)
