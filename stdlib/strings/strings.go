package strings

type Builder struct {
	buf []byte
}

func Contains(s string, substr string) bool {
	return Index(s, substr) >= 0
}

func HasPrefix(s string, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}

	for index := 0; index < len(prefix); index++ {
		if s[index] != prefix[index] {
			return false
		}
	}

	return true
}

func HasSuffix(s string, suffix string) bool {
	if len(suffix) > len(s) {
		return false
	}

	start := len(s) - len(suffix)
	for index := 0; index < len(suffix); index++ {
		if s[start+index] != suffix[index] {
			return false
		}
	}

	return true
}

func Index(s string, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

	limit := len(s) - len(substr)
	for index := 0; index <= limit; index++ {
		if HasPrefix(s[index:], substr) {
			return index
		}
	}

	return -1
}

func LastIndex(s string, substr string) int {
	if len(substr) == 0 {
		return len(s)
	}
	if len(substr) > len(s) {
		return -1
	}

	for index := len(s) - len(substr); index >= 0; index-- {
		if HasPrefix(s[index:], substr) {
			return index
		}
	}

	return -1
}

func Join(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}

	joined := elems[0]
	for index := 1; index < len(elems); index++ {
		joined += sep + elems[index]
	}

	return joined
}

func Cut(s string, sep string) (before string, after string, found bool) {
	if len(sep) == 0 {
		return "", s, true
	}

	index := Index(s, sep)
	if index < 0 {
		return s, "", false
	}

	return s[:index], s[index+len(sep):], true
}

func TrimPrefix(s string, prefix string) string {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}

	return s
}

func TrimSuffix(s string, suffix string) string {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}

	return s
}

func (builder *Builder) String() string {
	if builder == nil {
		return ""
	}

	return string(builder.buf)
}

func (builder *Builder) Len() int {
	if builder == nil {
		return 0
	}

	return len(builder.buf)
}

func (builder *Builder) Cap() int {
	if builder == nil {
		return 0
	}

	return cap(builder.buf)
}

func (builder *Builder) Reset() {
	if builder == nil {
		return
	}

	builder.buf = builder.buf[:0]
}

func (builder *Builder) Grow(n int) {
	if builder == nil || n <= 0 {
		return
	}
	if cap(builder.buf)-len(builder.buf) >= n {
		return
	}

	grown := make([]byte, len(builder.buf), len(builder.buf)+n)
	copy(grown, builder.buf)
	builder.buf = grown
}

func (builder *Builder) Write(data []byte) (int, error) {
	if builder == nil {
		return 0, nil
	}

	builder.buf = append(builder.buf, data...)
	return len(data), nil
}

func (builder *Builder) WriteByte(value byte) error {
	if builder == nil {
		return nil
	}

	builder.buf = append(builder.buf, value)
	return nil
}

func (builder *Builder) WriteString(value string) (int, error) {
	if builder == nil {
		return 0, nil
	}

	builder.buf = append(builder.buf, value...)
	return len(value), nil
}
