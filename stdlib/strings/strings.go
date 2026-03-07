package strings

import "io"

type Builder struct {
	buf []byte
}

type Reader struct {
	value string
	index int64
	prev  int64
}

type readerError struct {
	text string
}

func (err *readerError) Error() string {
	return err.text
}

func NewReader(s string) *Reader {
	return &Reader{
		value: s,
		prev:  -1,
	}
}

func (reader *Reader) Len() int {
	if reader == nil || reader.index >= int64(len(reader.value)) {
		return 0
	}

	return len(reader.value) - int(reader.index)
}

func (reader *Reader) Size() int64 {
	if reader == nil {
		return 0
	}

	return int64(len(reader.value))
}

func (reader *Reader) Reset(s string) {
	if reader == nil {
		return
	}

	reader.value = s
	reader.index = 0
	reader.prev = -1
}

func (reader *Reader) Read(p []byte) (n int, err error) {
	if reader == nil {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	if reader.index >= int64(len(reader.value)) {
		reader.prev = -1
		return 0, io.EOF
	}

	start := int(reader.index)
	n = copy(p, reader.value[start:])
	reader.index += int64(n)
	reader.prev = -1
	return n, nil
}

func (reader *Reader) ReadAt(p []byte, off int64) (n int, err error) {
	if reader == nil {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	if off < 0 {
		return 0, &readerError{text: "strings.Reader.ReadAt: negative offset"}
	}
	if off >= int64(len(reader.value)) {
		return 0, io.EOF
	}

	n = copy(p, reader.value[int(off):])
	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}

func (reader *Reader) ReadByte() (byte, error) {
	if reader == nil || reader.index >= int64(len(reader.value)) {
		return 0, io.EOF
	}

	reader.prev = reader.index
	value := reader.value[int(reader.index)]
	reader.index++
	return value, nil
}

func (reader *Reader) UnreadByte() error {
	if reader == nil || reader.prev < 0 {
		return &readerError{text: "strings.Reader.UnreadByte: no byte to unread"}
	}

	reader.index = reader.prev
	reader.prev = -1
	return nil
}

func (reader *Reader) Seek(offset int64, whence int) (int64, error) {
	if reader == nil {
		return 0, &readerError{text: "strings.Reader.Seek: nil reader"}
	}

	base := int64(0)
	switch whence {
	case io.SeekStart:
		base = 0
	case io.SeekCurrent:
		base = reader.index
	case io.SeekEnd:
		base = int64(len(reader.value))
	default:
		return reader.index, &readerError{text: "strings.Reader.Seek: invalid whence"}
	}

	position := base + offset
	if position < 0 {
		return reader.index, &readerError{text: "strings.Reader.Seek: negative position"}
	}

	reader.index = position
	reader.prev = -1
	return reader.index, nil
}

func (reader *Reader) WriteTo(writer io.Writer) (n int64, err error) {
	if reader == nil || reader.index >= int64(len(reader.value)) {
		return 0, nil
	}

	start := int(reader.index)
	wrote, err := writer.Write([]byte(reader.value[start:]))
	reader.index += int64(wrote)
	reader.prev = -1
	n = int64(wrote)
	if err != nil {
		return n, err
	}
	if wrote != len(reader.value)-start {
		return n, io.ErrShortWrite
	}

	return n, nil
}

func Contains(s string, substr string) bool {
	return Index(s, substr) >= 0
}

func Split(s string, sep string) []string {
	return SplitN(s, sep, -1)
}

func SplitN(s string, sep string, n int) []string {
	if n == 0 {
		return nil
	}
	if sep == "" {
		return splitRunesASCII(s, n)
	}
	if n == 1 {
		return []string{s}
	}

	var parts []string
	start := 0
	for start <= len(s) {
		if n > 0 && len(parts)+1 >= n {
			parts = append(parts, s[start:])
			return parts
		}

		index := Index(s[start:], sep)
		if index < 0 {
			parts = append(parts, s[start:])
			return parts
		}

		index += start
		parts = append(parts, s[start:index])
		start = index + len(sep)
	}

	return append(parts, "")
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

func TrimSpace(s string) string {
	start := 0
	for start < len(s) && isASCIISpace(s[start]) {
		start++
	}

	end := len(s)
	for end > start && isASCIISpace(s[end-1]) {
		end--
	}

	return s[start:end]
}

func Fields(s string) []string {
	var fields []string
	start := -1

	for index := 0; index < len(s); index++ {
		if isASCIISpace(s[index]) {
			if start >= 0 {
				fields = append(fields, s[start:index])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = index
		}
	}
	if start >= 0 {
		fields = append(fields, s[start:])
	}

	return fields
}

func ReplaceAll(s string, old string, new string) string {
	if old == "" {
		if len(s) == 0 {
			return new
		}

		replaced := new
		for index := 0; index < len(s); index++ {
			replaced += s[index:index+1] + new
		}
		return replaced
	}

	index := Index(s, old)
	if index < 0 {
		return s
	}

	replaced := ""
	start := 0
	for index >= 0 {
		replaced += s[start:index] + new
		start = index + len(old)
		next := Index(s[start:], old)
		if next < 0 {
			break
		}
		index = start + next
	}

	return replaced + s[start:]
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

func splitRunesASCII(s string, n int) []string {
	if len(s) == 0 {
		return []string{}
	}
	if n > 0 && n == 1 {
		return []string{s}
	}

	parts := make([]string, 0, len(s))
	for index := 0; index < len(s); index++ {
		if n > 0 && len(parts)+1 >= n {
			parts = append(parts, s[index:])
			return parts
		}
		parts = append(parts, s[index:index+1])
	}

	return parts
}

func isASCIISpace(value byte) bool {
	switch value {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	}

	return false
}
