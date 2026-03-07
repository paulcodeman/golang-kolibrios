package fmt

import (
	"errors"
	"io"
	"os"
)

type Stringer interface {
	String() string
}

type buffer struct {
	data []byte
}

func (buffer *buffer) Write(data []byte) (int, error) {
	buffer.data = append(buffer.data, data...)
	return len(data), nil
}

func (buffer *buffer) String() string {
	return string(buffer.data)
}

func (buffer *buffer) writeString(value string) {
	if value == "" {
		return
	}

	buffer.data = append(buffer.data, value...)
}

func (buffer *buffer) writeByte(value byte) {
	buffer.data = append(buffer.data, value)
}

func writeRendered(writer io.Writer, text string) (n int, err error) {
	if text == "" {
		return 0, nil
	}

	written, err := io.WriteString(writer, text)
	if err != nil {
		return written, err
	}
	if written != len(text) {
		return written, io.ErrShortWrite
	}

	return written, nil
}

func renderPrint(values ...interface{}) string {
	buffer := &buffer{}
	for index := 0; index < len(values); index++ {
		buffer.writeString(formatValue(values[index], 'v'))
	}

	return buffer.String()
}

func renderPrintln(values ...interface{}) string {
	buffer := &buffer{}
	for index := 0; index < len(values); index++ {
		if index > 0 {
			buffer.writeByte(' ')
		}
		buffer.writeString(formatValue(values[index], 'v'))
	}
	buffer.writeByte('\n')

	return buffer.String()
}

func renderPrintf(format string, values ...interface{}) string {
	buffer := &buffer{}
	valueIndex := 0
	textStart := 0

	for index := 0; index < len(format); index++ {
		if format[index] != '%' {
			continue
		}

		buffer.writeString(format[textStart:index])
		textStart = index + 1

		if textStart >= len(format) {
			buffer.writeString("%!(NOVERB)")
			return buffer.String()
		}

		verb := format[textStart]
		textStart++
		if verb == '%' {
			buffer.writeByte('%')
			index = textStart - 1
			continue
		}

		if valueIndex >= len(values) {
			buffer.writeString(missingVerb(verb))
			index = textStart - 1
			continue
		}

		buffer.writeString(formatValue(values[valueIndex], verb))
		valueIndex++
		index = textStart - 1
	}

	if textStart < len(format) {
		buffer.writeString(format[textStart:])
	}

	return buffer.String()
}

func Sprint(values ...interface{}) string {
	buffer := &buffer{}
	_, _ = Fprint(buffer, values...)
	return buffer.String()
}

func Sprintln(values ...interface{}) string {
	buffer := &buffer{}
	_, _ = Fprintln(buffer, values...)
	return buffer.String()
}

func Sprintf(format string, values ...interface{}) string {
	buffer := &buffer{}
	_, _ = Fprintf(buffer, format, values...)
	return buffer.String()
}

func Fprint(writer io.Writer, values ...interface{}) (n int, err error) {
	return writeRendered(writer, renderPrint(values...))
}

func Fprintln(writer io.Writer, values ...interface{}) (n int, err error) {
	return writeRendered(writer, renderPrintln(values...))
}

func Fprintf(writer io.Writer, format string, values ...interface{}) (n int, err error) {
	return writeRendered(writer, renderPrintf(format, values...))
}

func Print(values ...interface{}) (n int, err error) {
	return Fprint(os.DefaultStdout(), values...)
}

func Println(values ...interface{}) (n int, err error) {
	return Fprintln(os.DefaultStdout(), values...)
}

func Printf(format string, values ...interface{}) (n int, err error) {
	return Fprintf(os.DefaultStdout(), format, values...)
}

func Fscan(reader io.Reader, values ...interface{}) (n int, err error) {
	return scanValues(reader, false, values...)
}

func Fscanln(reader io.Reader, values ...interface{}) (n int, err error) {
	return scanValues(reader, true, values...)
}

func Scan(values ...interface{}) (n int, err error) {
	return Fscan(os.DefaultStdin(), values...)
}

func Scanln(values ...interface{}) (n int, err error) {
	return Fscanln(os.DefaultStdin(), values...)
}

func Errorf(format string, values ...interface{}) error {
	return errors.New(Sprintf(format, values...))
}

var errScanSyntax = errors.New("invalid scan syntax")
var errScanTarget = errors.New("unsupported scan target")
var errScanNewline = errors.New("unexpected newline")
var errScanTrailing = errors.New("expected newline")

type scanReader struct {
	reader     io.Reader
	byteBuffer [1]byte
	pending    byte
	haveByte   bool
	lineEnded  bool
}

func scanValues(reader io.Reader, lineMode bool, values ...interface{}) (n int, err error) {
	scanner := &scanReader{reader: reader}

	for index := 0; index < len(values); index++ {
		token, readErr := scanner.readToken(lineMode)
		if readErr != nil {
			return n, readErr
		}
		if assignErr := scanAssign(values[index], token); assignErr != nil {
			return n, assignErr
		}
		n++
	}

	if lineMode {
		if err = scanner.consumeLineTail(); err != nil {
			return n, err
		}
	}

	return n, nil
}

func (scanner *scanReader) readToken(lineMode bool) (string, error) {
	if lineMode && scanner.lineEnded {
		return "", errScanNewline
	}

	for {
		value, err := scanner.readByte()
		if err != nil {
			return "", err
		}
		if isScanNewline(value) {
			if lineMode {
				scanner.lineEnded = true
				return "", errScanNewline
			}
			continue
		}
		if isScanHorizontalSpace(value) {
			continue
		}

		scanner.unreadByte(value)
		break
	}

	token := &buffer{}
	for {
		value, err := scanner.readByte()
		if err != nil {
			if len(token.data) > 0 {
				return token.String(), nil
			}
			return "", err
		}
		if isScanNewline(value) {
			if lineMode {
				scanner.lineEnded = true
			}
			break
		}
		if isScanHorizontalSpace(value) {
			break
		}

		token.writeByte(value)
	}

	if len(token.data) == 0 {
		return "", io.EOF
	}

	return token.String(), nil
}

func (scanner *scanReader) consumeLineTail() error {
	if scanner.lineEnded {
		return nil
	}

	for {
		value, err := scanner.readByte()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if isScanNewline(value) {
			scanner.lineEnded = true
			return nil
		}
		if isScanHorizontalSpace(value) {
			continue
		}

		return errScanTrailing
	}
}

func (scanner *scanReader) readByte() (byte, error) {
	if scanner.haveByte {
		scanner.haveByte = false
		return scanner.pending, nil
	}

	for {
		read, err := scanner.reader.Read(scanner.byteBuffer[:])
		if read > 0 {
			return scanner.byteBuffer[0], nil
		}
		if err != nil {
			return 0, err
		}
	}
}

func (scanner *scanReader) unreadByte(value byte) {
	scanner.pending = value
	scanner.haveByte = true
}

func isScanHorizontalSpace(value byte) bool {
	switch value {
	case ' ', '\t', '\v', '\f':
		return true
	}

	return false
}

func isScanNewline(value byte) bool {
	return value == '\n' || value == '\r'
}

func scanAssign(target interface{}, token string) error {
	switch typed := target.(type) {
	case *string:
		if typed == nil {
			return errScanTarget
		}
		*typed = token
		return nil
	case *bool:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseBoolToken(token)
		if err != nil {
			return err
		}
		*typed = value
		return nil
	case *int:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseSignedToken(token, intBitSize())
		if err != nil {
			return err
		}
		*typed = int(value)
		return nil
	case *int8:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseSignedToken(token, 8)
		if err != nil {
			return err
		}
		*typed = int8(value)
		return nil
	case *int16:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseSignedToken(token, 16)
		if err != nil {
			return err
		}
		*typed = int16(value)
		return nil
	case *int32:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseSignedToken(token, 32)
		if err != nil {
			return err
		}
		*typed = int32(value)
		return nil
	case *int64:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseSignedToken(token, 64)
		if err != nil {
			return err
		}
		*typed = value
		return nil
	case *uint:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseUnsignedToken(token, uintBitSize())
		if err != nil {
			return err
		}
		*typed = uint(value)
		return nil
	case *uint8:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseUnsignedToken(token, 8)
		if err != nil {
			return err
		}
		*typed = uint8(value)
		return nil
	case *uint16:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseUnsignedToken(token, 16)
		if err != nil {
			return err
		}
		*typed = uint16(value)
		return nil
	case *uint32:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseUnsignedToken(token, 32)
		if err != nil {
			return err
		}
		*typed = uint32(value)
		return nil
	case *uint64:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseUnsignedToken(token, 64)
		if err != nil {
			return err
		}
		*typed = value
		return nil
	case *uintptr:
		if typed == nil {
			return errScanTarget
		}
		value, err := parseUnsignedToken(token, uintBitSize())
		if err != nil {
			return err
		}
		*typed = uintptr(value)
		return nil
	}

	return errScanTarget
}

func parseBoolToken(token string) (bool, error) {
	if equalFoldASCII(token, "true") || equalFoldASCII(token, "t") || token == "1" {
		return true, nil
	}
	if equalFoldASCII(token, "false") || equalFoldASCII(token, "f") || token == "0" {
		return false, nil
	}

	return false, errScanSyntax
}

func parseSignedToken(token string, bits uint) (int64, error) {
	if token == "" {
		return 0, errScanSyntax
	}

	negative := false
	switch token[0] {
	case '+':
		token = token[1:]
	case '-':
		negative = true
		token = token[1:]
	}
	if token == "" {
		return 0, errScanSyntax
	}

	base := uint64(10)
	if len(token) > 2 && token[0] == '0' && (token[1] == 'x' || token[1] == 'X') {
		base = 16
		token = token[2:]
	}
	if token == "" {
		return 0, errScanSyntax
	}

	limit := maxSignedMagnitude(bits, negative)

	value, err := parseUnsignedWithBase(token, base, limit)
	if err != nil {
		return 0, err
	}
	if negative {
		if value == uint64(1)<<(bits-1) {
			return -int64(value), nil
		}
		return -int64(value), nil
	}

	return int64(value), nil
}

func parseUnsignedToken(token string, bits uint) (uint64, error) {
	if token == "" {
		return 0, errScanSyntax
	}
	if token[0] == '+' {
		token = token[1:]
	}
	if token == "" || token[0] == '-' {
		return 0, errScanSyntax
	}

	base := uint64(10)
	if len(token) > 2 && token[0] == '0' && (token[1] == 'x' || token[1] == 'X') {
		base = 16
		token = token[2:]
	}
	if token == "" {
		return 0, errScanSyntax
	}

	return parseUnsignedWithBase(token, base, maxUnsignedForBits(bits))
}

func maxSignedMagnitude(bits uint, negative bool) uint64 {
	if bits >= 64 {
		if negative {
			return uint64(1) << 63
		}

		return (uint64(1) << 63) - 1
	}
	if negative {
		return uint64(1) << (bits - 1)
	}

	return (uint64(1) << (bits - 1)) - 1
}

func maxUnsignedForBits(bits uint) uint64 {
	if bits >= 64 {
		return ^uint64(0)
	}

	return (uint64(1) << bits) - 1
}

func parseUnsignedWithBase(token string, base uint64, limit uint64) (uint64, error) {
	value := uint64(0)

	for index := 0; index < len(token); index++ {
		digit, ok := digitValue(token[index])
		if !ok || digit >= base {
			return 0, errScanSyntax
		}

		next := value*base + digit
		if next < value || next > limit {
			return 0, errScanSyntax
		}
		value = next
	}

	return value, nil
}

func digitValue(value byte) (uint64, bool) {
	switch {
	case value >= '0' && value <= '9':
		return uint64(value - '0'), true
	case value >= 'a' && value <= 'f':
		return uint64(value-'a') + 10, true
	case value >= 'A' && value <= 'F':
		return uint64(value-'A') + 10, true
	}

	return 0, false
}

func equalFoldASCII(left string, right string) bool {
	if len(left) != len(right) {
		return false
	}

	for index := 0; index < len(left); index++ {
		if lowerASCII(left[index]) != lowerASCII(right[index]) {
			return false
		}
	}

	return true
}

func lowerASCII(value byte) byte {
	if value >= 'A' && value <= 'Z' {
		return value + ('a' - 'A')
	}

	return value
}

func intBitSize() uint {
	return uintBitSize()
}

func uintBitSize() uint {
	return uint(^uint(0)>>63)*32 + 32
}

var decimalDigits = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var lowerHexDigits = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
var upperHexDigits = [...]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
var decimalPowers = [...]uint64{
	10000000000000000000,
	1000000000000000000,
	100000000000000000,
	10000000000000000,
	1000000000000000,
	100000000000000,
	10000000000000,
	1000000000000,
	100000000000,
	10000000000,
	1000000000,
	100000000,
	10000000,
	1000000,
	100000,
	10000,
	1000,
	100,
	10,
	1,
}

func formatValue(value interface{}, verb byte) string {
	if verb == 'w' {
		verb = 'v'
	}

	if value == nil {
		if verb == 'v' || verb == 's' {
			return "<nil>"
		}

		return unsupportedVerb(verb)
	}

	switch typed := value.(type) {
	case string:
		return formatStringValue(typed, verb)
	case []byte:
		return formatBytesValue(typed, verb)
	case bool:
		if verb == 't' || verb == 'v' {
			return formatBool(typed)
		}
	case int:
		return formatSignedValue(int64(typed), verb)
	case int8:
		return formatSignedValue(int64(typed), verb)
	case int16:
		return formatSignedValue(int64(typed), verb)
	case int32:
		return formatSignedValue(int64(typed), verb)
	case int64:
		return formatSignedValue(typed, verb)
	case uint:
		return formatUnsignedValue(uint64(typed), verb)
	case uint8:
		return formatUnsignedValue(uint64(typed), verb)
	case uint16:
		return formatUnsignedValue(uint64(typed), verb)
	case uint32:
		return formatUnsignedValue(uint64(typed), verb)
	case uint64:
		return formatUnsignedValue(typed, verb)
	case uintptr:
		return formatUnsignedValue(uint64(typed), verb)
	}

	if err, ok := value.(error); ok {
		if verb == 'v' || verb == 's' {
			return err.Error()
		}
	}

	if stringer, ok := value.(Stringer); ok {
		if verb == 'v' || verb == 's' {
			return stringer.String()
		}
	}

	return unsupportedVerb(verb)
}

func formatStringValue(value string, verb byte) string {
	switch verb {
	case 's', 'v':
		return value
	case 'x':
		return formatHexBytes([]byte(value), false)
	case 'X':
		return formatHexBytes([]byte(value), true)
	}

	return unsupportedVerb(verb)
}

func formatBytesValue(value []byte, verb byte) string {
	switch verb {
	case 's', 'v':
		return string(value)
	case 'x':
		return formatHexBytes(value, false)
	case 'X':
		return formatHexBytes(value, true)
	}

	return unsupportedVerb(verb)
}

func formatSignedValue(value int64, verb byte) string {
	switch verb {
	case 'd', 'v':
		return formatInt64Decimal(value)
	case 'x':
		return formatUint64Hex(uint64(value), lowerHexDigits[:])
	case 'X':
		return formatUint64Hex(uint64(value), upperHexDigits[:])
	case 'c':
		return string([]byte{byte(value)})
	}

	return unsupportedVerb(verb)
}

func formatUnsignedValue(value uint64, verb byte) string {
	switch verb {
	case 'd', 'v':
		return formatUint64Decimal(value)
	case 'x':
		return formatUint64Hex(value, lowerHexDigits[:])
	case 'X':
		return formatUint64Hex(value, upperHexDigits[:])
	case 'c':
		return string([]byte{byte(value)})
	}

	return unsupportedVerb(verb)
}

func formatBool(value bool) string {
	if value {
		return "true"
	}

	return "false"
}

func formatInt64Decimal(value int64) string {
	if value < 0 {
		return "-" + formatUint64Decimal(uint64(^value)+1)
	}

	return formatUint64Decimal(uint64(value))
}

func formatUint64Decimal(value uint64) string {
	if value == 0 {
		return "0"
	}

	text := ""
	started := false

	for index := 0; index < len(decimalPowers); index++ {
		digit := uint32(0)
		for value >= decimalPowers[index] {
			value -= decimalPowers[index]
			digit++
		}

		if digit != 0 || started {
			text += decimalDigits[digit]
			started = true
		}
	}

	return text
}

func formatUint64Hex(value uint64, digits []string) string {
	if value == 0 {
		return "0"
	}

	text := ""
	started := false

	for shift := uint(60); ; shift -= 4 {
		digit := uint32((value >> shift) & 0x0F)
		if digit != 0 || started {
			text += digits[digit]
			started = true
		}

		if shift == 0 {
			break
		}
	}

	return text
}

func formatHexBytes(value []byte, upper bool) string {
	if len(value) == 0 {
		return ""
	}

	digits := lowerHexDigits[:]
	if upper {
		digits = upperHexDigits[:]
	}

	text := ""
	for index := 0; index < len(value); index++ {
		text += digits[uint32(value[index]>>4)]
		text += digits[uint32(value[index]&0x0F)]
	}

	return text
}

func missingVerb(verb byte) string {
	return "%!" + string([]byte{verb}) + "(MISSING)"
}

func unsupportedVerb(verb byte) string {
	return "%!" + string([]byte{verb}) + "(UNSUPPORTED)"
}
