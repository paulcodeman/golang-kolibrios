package runtimeprobe

func EmptyInterfaceEq() bool {
	var left interface{} = "same"
	var right interface{} = "same"

	return left == right
}
