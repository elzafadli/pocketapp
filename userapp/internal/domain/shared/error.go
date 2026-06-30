package shared

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Code    string                  `json:"code"`
	Message string                  `json:"message"`
	Details []ValidationErrorDetail `json:"details"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
