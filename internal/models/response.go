package models

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type DataResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewErrorResponse(code int, error string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   error,
		Code:    code,
	}
}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Data:    data,
	}
}

func NewDataResponse(message string, data interface{}) DataResponse {
	return DataResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}