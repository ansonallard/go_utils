package ierr

type BadRequestError struct {
	message string
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{message: message}
}

func (e *BadRequestError) Error() string {
	if e.message != "" {
		return e.message
	}
	return "BadRequestError"
}
