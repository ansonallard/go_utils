package ierr

type ForbiddenError struct{}

func (e *ForbiddenError) Error() string {
	return "ForbiddenError"
}
