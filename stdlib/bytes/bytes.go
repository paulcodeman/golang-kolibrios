package bytes

import "io"

type Buffer struct {
	buf []byte
}

type Reader struct {
	value []byte
	index int64
	prev  int64
}

type readerError struct {
	text string
}

func (err *readerError) Error() string {
	return err.text
}

func NewReader(data []byte) *Reader {
	return &Reader{
		value: data,
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

func (reader *Reader) Reset(data []byte) {
	if reader == nil {
		return
	}

	reader.value = data
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
		return 0, &readerError{text: "bytes.Reader.ReadAt: negative offset"}
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
		return &readerError{text: "bytes.Reader.UnreadByte: no byte to unread"}
	}

	reader.index = reader.prev
	reader.prev = -1
	return nil
}

func (reader *Reader) Seek(offset int64, whence int) (int64, error) {
	if reader == nil {
		return 0, &readerError{text: "bytes.Reader.Seek: nil reader"}
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
		return reader.index, &readerError{text: "bytes.Reader.Seek: invalid whence"}
	}

	position := base + offset
	if position < 0 {
		return reader.index, &readerError{text: "bytes.Reader.Seek: negative position"}
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
	wrote, err := writer.Write(reader.value[start:])
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

func Split(s []byte, sep []byte) [][]byte {
	return SplitN(s, sep, -1)
}

func SplitN(s []byte, sep []byte, n int) [][]byte {
	if n == 0 {
		return nil
	}
	if len(sep) == 0 {
		return splitBytes(s, n)
	}
	if n == 1 {
		return [][]byte{s}
	}

	var parts [][]byte
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

	return append(parts, []byte{})
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

func TrimSpace(s []byte) []byte {
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

func Fields(s []byte) [][]byte {
	var fields [][]byte
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

func ReplaceAll(s []byte, old []byte, new []byte) []byte {
	if len(old) == 0 {
		if len(s) == 0 {
			return append([]byte{}, new...)
		}

		replaced := make([]byte, 0, len(s)+(len(s)+1)*len(new))
		replaced = append(replaced, new...)
		for index := 0; index < len(s); index++ {
			replaced = append(replaced, s[index])
			replaced = append(replaced, new...)
		}
		return replaced
	}

	index := Index(s, old)
	if index < 0 {
		return append([]byte{}, s...)
	}

	replaced := make([]byte, 0, len(s))
	start := 0
	for index >= 0 {
		replaced = append(replaced, s[start:index]...)
		replaced = append(replaced, new...)
		start = index + len(old)
		next := Index(s[start:], old)
		if next < 0 {
			break
		}
		index = start + next
	}
	replaced = append(replaced, s[start:]...)
	return replaced
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

func splitBytes(s []byte, n int) [][]byte {
	if len(s) == 0 {
		return [][]byte{}
	}
	if n > 0 && n == 1 {
		return [][]byte{s}
	}

	parts := make([][]byte, 0, len(s))
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
