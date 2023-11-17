package errors

type IErrorMessage interface {
	Message() string
	HttpStatus() int
}

type ErrorMessage struct {
	ErrorMessageText string `json:"message"`
	ErrorStatus      int    `json:"status"`
}

func (e *ErrorMessage) Message() string {
	return e.ErrorMessageText
}

func (e *ErrorMessage) HttpStatus() int {
	return e.ErrorStatus
}

func NewErrorMessage(message string, httpStatusCode int) IErrorMessage {
	return &ErrorMessage{ErrorMessageText: message, ErrorStatus: httpStatusCode}
}
