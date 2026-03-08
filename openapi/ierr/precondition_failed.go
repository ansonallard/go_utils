package ierr

type PreConditionFailed struct {
	message string
}

func NewPreConditionFailed(message string) *PreConditionFailed {
	return &PreConditionFailed{message: message}
}

func (e *PreConditionFailed) Error() string {
	if e.message != "" {
		return e.message
	}
	return "PreConditionFailed"
}
