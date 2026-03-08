package ierr

type TooManyRequestsError struct {
	message string
}

func NewTooManyRequestsError(message string) *TooManyRequestsError {
	return &TooManyRequestsError{message: message}
}

func (e *TooManyRequestsError) Error() string {
	if e.message != "" {
		return e.message
	}
	return "TooManyRequestsError"
}
