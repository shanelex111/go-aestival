package error_code

import (
	"net/http"

	"github.com/shanelex111/go-common/pkg/response"
)

var (
	AuthBadRequest = &response.Error{
		Status:  http.StatusBadRequest,
		Code:    101400,
		Message: http.StatusText(http.StatusBadRequest),
	}
	AuthUnauthorized = &response.Error{
		Status:  http.StatusUnauthorized,
		Code:    101401,
		Message: http.StatusText(http.StatusUnauthorized),
	}
	AuthInvalidPassword = &response.Error{
		Status:  http.StatusUnauthorized,
		Code:    101402,
		Message: "Invalid consumer or password.",
	}
	AuthForbidden = &response.Error{
		Status:  http.StatusForbidden,
		Code:    101403,
		Message: http.StatusText(http.StatusForbidden),
	}
	AuthNotFound = &response.Error{
		Status:  http.StatusNotFound,
		Code:    101404,
		Message: http.StatusText(http.StatusNotFound),
	}

	AuthInvalidVerificationCode = &response.Error{
		Status:  http.StatusUnauthorized,
		Code:    101405,
		Message: "Invalid verification code.",
	}

	AuthInternalServerError = &response.Error{
		Status:  http.StatusInternalServerError,
		Code:    101500,
		Message: http.StatusText(http.StatusInternalServerError),
	}
)
