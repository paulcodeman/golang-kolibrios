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

func Errorf(format string, values ...interface{}) error {
	return errors.New(Sprintf(format, values...))
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
