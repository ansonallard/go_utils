package ierr

type ConflictError struct {
	message string
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{message: message}
}

func (e *ConflictError) Error() string {
	if e.message != "" {
		return e.message
	}
	return "ConflictError"
}
