package runtimeprobe

type labeler interface {
	Label() string
}

type box struct {
	value string
}

func (b box) Label() string {
	return b.value
}

func AssertIface(v interface{}) labeler {
	return v.(labeler)
}

func AssertIfaceOK(v interface{}) (labeler, bool) {
	out, ok := v.(labeler)
	return out, ok
}
