package handler

import (
	"errors"
	"net/http"
	"web_app/domain"
)

func httpStatusFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrUsernameExist):
		return http.StatusConflict // 409
	case errors.Is(err, domain.ErrWrongPassword):
		return http.StatusUnauthorized // 401
	default:
		return http.StatusInternalServerError // 500
	}
}
