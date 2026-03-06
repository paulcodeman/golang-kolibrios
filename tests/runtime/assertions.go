package runtimeprobe

func AssertString(v interface{}) string {
	return v.(string)
}

func AssertStringOK(v interface{}) (string, bool) {
	s, ok := v.(string)
	return s, ok
}
