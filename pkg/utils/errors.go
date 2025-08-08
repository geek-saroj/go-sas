package utils

import (
	"net/http"
	"sas-pro/internal/models"
)

func UnauthorizedError() models.ErrorResponse {
	return models.NewErrorResponse(
		http.StatusUnauthorized,
		"Authorization header required or invalid",
	)
}

func ForbiddenError() models.ErrorResponse {
	return models.NewErrorResponse(
		http.StatusForbidden,
		"You don't have permission to access this resource",
	)
}

func BadRequestError(msg ...string) models.ErrorResponse {
	return models.NewErrorResponse(
		http.StatusBadRequest,
		fallbackMessage(getFirst(msg), "Invalid request"),
	)
}

func getFirst(msg []string) string {
	if len(msg) > 0 {
		return msg[0]
	}
	return ""
}

func fallbackMessage(custom, fallback string) string {
	if custom != "" {
		return custom
	}
	return fallback
}
