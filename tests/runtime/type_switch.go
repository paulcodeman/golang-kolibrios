package runtimeprobe

func TypeSwitch(v interface{}) int {
	switch v.(type) {
	case string:
		return 1
	case int:
		return 2
	default:
		return 0
	}
}
