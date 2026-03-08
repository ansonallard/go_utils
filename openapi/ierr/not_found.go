package ierr

type NotFoundError struct {
	message string
}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{message: message}
}

func (e *NotFoundError) Error() string {
	if e.message != "" {
		return e.message
	}
	return "NotFound"
}
