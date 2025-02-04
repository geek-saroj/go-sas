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