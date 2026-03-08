package ierr

type UnAuthorizedError struct {
	message string
}

func NewUnAuthorizedError(message string) *UnAuthorizedError {
	return &UnAuthorizedError{message: message}
}

func (e *UnAuthorizedError) Error() string {
	if e.message != "" {
		return e.message
	}
	return "Unauthorized"
}
