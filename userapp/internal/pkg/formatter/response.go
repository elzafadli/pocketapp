package formatter

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	TraceID    string `json:"traceId,omitempty"`
	StatusCode *int   `json:"statusCodeClient,omitempty"`
	ErrorList  any    `json:"errorList,omitempty"`
}

type ValidationErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewSuccessResponse(status Status, data any) *Response {
	return &Response{
		Status:  status.String(),
		Message: "success",
		Data:    data,
	}
}

func NewErrorResponse(status Status, message string, id string) *Response {
	return &Response{
		Status:  status.String(),
		Message: message,
		TraceID: id,
	}
}

func NewErrorResponseList(status Status, message string, id string, err any) *Response {
	return &Response{
		Status:    status.String(),
		Message:   message,
		TraceID:   id,
		ErrorList: err,
	}
}

// NewBadRequestResponse returns a 400-style error body.
func NewBadRequestResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    "BAD_REQUEST",
		Message: message,
	}
}

// NewNotFoundResponse returns a 404-style error body.
func NewNotFoundResponse(code, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// NewInternalErrorResponse returns a 500-style error body.
func NewInternalErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: message,
	}
}

// NewValidationErrorResponse returns a 400-style validation error body with per-field details.
func NewValidationErrorResponse(details any) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		Code:    "VALIDATION_ERROR",
		Message: "Validation error",
		Details: details,
	}
}
