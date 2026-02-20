package ierr

type Forbidden struct{}

func (e *Forbidden) Error() string {
	return "Forbidden"
}
