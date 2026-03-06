package runtimeprobe

type source interface {
	Label() string
}

type target interface {
	Label() string
}

func CastIface(v source) target {
	return v.(target)
}

func CastIfaceOK(v source) (target, bool) {
	out, ok := v.(target)
	return out, ok
}
