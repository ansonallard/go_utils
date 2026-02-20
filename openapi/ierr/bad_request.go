package ierr

type BadRequestError struct{}

func (e *BadRequestError) Error() string {
	return "BadRequestError"
}
