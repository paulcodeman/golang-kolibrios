package runtimeprobe

type unwrapError interface {
	Unwrap() error
}

type wrappedError struct {
	cause error
}

func (err wrappedError) Error() string {
	return "wrapped"
}

func (err wrappedError) Unwrap() error {
	return err.cause
}

func CanUnwrap(err error) bool {
	var any interface{} = err

	_, ok := any.(unwrapError)
	return ok
}
