package formatter

type Response struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	TraceID    string `json:"traceId,omitempty"`
	StatusCode *int   `json:"statusCodeClient,omitempty"`
	ErrorList  any    `json:"errorList,omitempty"`
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
