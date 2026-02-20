package ierr

type PreConditionFailed struct{}

func (e *PreConditionFailed) Error() string {
	return "PreConditionFailed"
}
