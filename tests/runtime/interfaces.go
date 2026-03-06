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

func InterfaceEq() bool {
	var left labeler = box{value: "same"}
	var right labeler = box{value: "same"}

	return left.Label() == right.Label() && left == right
}
