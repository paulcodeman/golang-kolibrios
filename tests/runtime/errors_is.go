package runtimeprobe

type matcher interface {
	Is(error) bool
}

type matchError struct {
	cause error
}

func (err matchError) Error() string {
	return "match"
}

func (err matchError) Is(target error) bool {
	return target == err.cause
}

func HasCustomIs(err error, target error) bool {
	var any interface{} = err

	value, ok := any.(matcher)
	if !ok {
		return false
	}

	return value.Is(target)
}
