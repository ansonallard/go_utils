package ierr

type ForbiddenError struct {
	message string
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{message: message}
}

func (e *ForbiddenError) Error() string {
	if e.message != "" {
		return e.message
	}
	return "ForbiddenError"
}
