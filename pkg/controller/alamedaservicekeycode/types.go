package alamedaservicekeycode

type reconcileError struct {
	err error
}

func (e *reconcileError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}
