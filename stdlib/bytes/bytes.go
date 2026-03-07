package bytes

type Buffer struct {
	buf []byte
}

func Contains(s []byte, subslice []byte) bool {
	return Index(s, subslice) >= 0
}

func Equal(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for index := 0; index < len(a); index++ {
		if a[index] != b[index] {
			return false
		}
	}

	return true
}

func HasPrefix(s []byte, prefix []byte) bool {
	return len(s) >= len(prefix) && Equal(s[:len(prefix)], prefix)
}

func HasSuffix(s []byte, suffix []byte) bool {
	return len(s) >= len(suffix) && Equal(s[len(s)-len(suffix):], suffix)
}

func Index(s []byte, sep []byte) int {
	if len(sep) == 0 {
		return 0
	}
	if len(sep) == 1 {
		return IndexByte(s, sep[0])
	}
	if len(sep) > len(s) {
		return -1
	}

	limit := len(s) - len(sep)
	for index := 0; index <= limit; index++ {
		if HasPrefix(s[index:], sep) {
			return index
		}
	}

	return -1
}

func IndexByte(s []byte, target byte) int {
	for index := 0; index < len(s); index++ {
		if s[index] == target {
			return index
		}
	}

	return -1
}

func Join(s [][]byte, sep []byte) []byte {
	if len(s) == 0 {
		return []byte{}
	}

	total := 0
	for index := 0; index < len(s); index++ {
		total += len(s[index])
	}
	if len(s) > 1 {
		total += len(sep) * (len(s) - 1)
	}

	joined := make([]byte, 0, total)
	for index := 0; index < len(s); index++ {
		if index > 0 {
			joined = append(joined, sep...)
		}
		joined = append(joined, s[index]...)
	}

	return joined
}

func Cut(s []byte, sep []byte) (before []byte, after []byte, found bool) {
	index := Index(s, sep)
	if index < 0 {
		return s, nil, false
	}

	return s[:index], s[index+len(sep):], true
}

func TrimPrefix(s []byte, prefix []byte) []byte {
	if HasPrefix(s, prefix) {
		return s[len(prefix):]
	}

	return s
}

func TrimSuffix(s []byte, suffix []byte) []byte {
	if HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}

	return s
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{buf: buf}
}

func NewBufferString(value string) *Buffer {
	buffer := &Buffer{}
	buffer.buf = append(buffer.buf, value...)
	return buffer
}

func (buffer *Buffer) Bytes() []byte {
	if buffer == nil {
		return nil
	}

	return buffer.buf
}

func (buffer *Buffer) String() string {
	if buffer == nil {
		return "<nil>"
	}

	return string(buffer.buf)
}

func (buffer *Buffer) Len() int {
	if buffer == nil {
		return 0
	}

	return len(buffer.buf)
}

func (buffer *Buffer) Cap() int {
	if buffer == nil {
		return 0
	}

	return cap(buffer.buf)
}

func (buffer *Buffer) Reset() {
	if buffer == nil {
		return
	}

	buffer.buf = buffer.buf[:0]
}

func (buffer *Buffer) Grow(n int) {
	if buffer == nil || n <= 0 {
		return
	}
	if cap(buffer.buf)-len(buffer.buf) >= n {
		return
	}

	grown := make([]byte, len(buffer.buf), len(buffer.buf)+n)
	copy(grown, buffer.buf)
	buffer.buf = grown
}

func (buffer *Buffer) Write(data []byte) (int, error) {
	if buffer == nil {
		return 0, nil
	}

	buffer.buf = append(buffer.buf, data...)
	return len(data), nil
}

func (buffer *Buffer) WriteByte(value byte) error {
	if buffer == nil {
		return nil
	}

	buffer.buf = append(buffer.buf, value)
	return nil
}

func (buffer *Buffer) WriteString(value string) (int, error) {
	if buffer == nil {
		return 0, nil
	}

	buffer.buf = append(buffer.buf, value...)
	return len(value), nil
}
