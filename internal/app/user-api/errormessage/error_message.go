package errormessage

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

func NewErrorMessage(message string, httpStatusCode int) ErrorMessage {
	return ErrorMessage{ErrorMessageText: message, ErrorStatus: httpStatusCode}
}
