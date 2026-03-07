package errors

type errorString struct {
	text string
}

func (err *errorString) Error() string {
	return err.text
}

type unwrapper interface {
	Unwrap() error
}

func New(text string) error {
	return &errorString{text: text}
}

func Unwrap(err error) error {
	if err == nil {
		return nil
	}

	var value interface{} = err
	unwrapped, ok := value.(unwrapper)
	if !ok {
		return nil
	}

	return unwrapped.Unwrap()
}

func Is(err error, target error) bool {
	if target == nil {
		return err == nil
	}

	for err != nil {
		if err == target {
			return true
		}

		err = Unwrap(err)
	}

	return false
}
