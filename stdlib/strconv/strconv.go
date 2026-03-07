package strconv

const IntSize = 32

type strconvError struct {
	text string
}

func (err *strconvError) Error() string {
	return err.text
}

var ErrRange = &strconvError{text: "value out of range"}
var ErrSyntax = &strconvError{text: "invalid syntax"}

type NumError struct {
	Func string
	Num  string
	Err  error
}

func (err *NumError) Error() string {
	if err == nil {
		return ""
	}
	if err.Err == nil {
		return err.Func + ": parsing " + quote(err.Num)
	}

	return err.Func + ": parsing " + quote(err.Num) + ": " + err.Err.Error()
}

func (err *NumError) Unwrap() error {
	if err == nil {
		return nil
	}

	return err.Err
}

func FormatBool(value bool) string {
	if value {
		return "true"
	}

	return "false"
}

func AppendBool(dst []byte, value bool) []byte {
	if value {
		return append(dst, 't', 'r', 'u', 'e')
	}

	return append(dst, 'f', 'a', 'l', 's', 'e')
}

func ParseBool(value string) (bool, error) {
	switch {
	case value == "1", value == "t", value == "T", equalFoldASCII(value, "true"):
		return true, nil
	case value == "0", value == "f", value == "F", equalFoldASCII(value, "false"):
		return false, nil
	}

	return false, &NumError{
		Func: "ParseBool",
		Num:  value,
		Err:  ErrSyntax,
	}
}

func Itoa(value int) string {
	return FormatInt(int64(value), 10)
}

func Atoi(value string) (int, error) {
	parsed, err := ParseInt(value, 10, 0)
	if err != nil {
		return 0, err
	}

	return int(parsed), nil
}

func FormatInt(value int64, base int) string {
	return string(AppendInt(nil, value, base))
}

func FormatUint(value uint64, base int) string {
	return string(AppendUint(nil, value, base))
}

func AppendInt(dst []byte, value int64, base int) []byte {
	normalizedBase := normalizeFormatBase(base)
	if value < 0 {
		dst = append(dst, '-')
		return appendUnsigned(dst, uint64(^value)+1, normalizedBase)
	}

	return appendUnsigned(dst, uint64(value), normalizedBase)
}

func AppendUint(dst []byte, value uint64, base int) []byte {
	return appendUnsigned(dst, value, normalizeFormatBase(base))
}

func ParseInt(value string, base int, bitSize int) (int64, error) {
	original := value
	if value == "" {
		return 0, numSyntax("ParseInt", original)
	}

	negative := false
	switch value[0] {
	case '+':
		value = value[1:]
	case '-':
		negative = true
		value = value[1:]
	}
	if value == "" {
		return 0, numSyntax("ParseInt", original)
	}

	bits := normalizeBitSize(bitSize)
	parseBase, digits, ok := normalizeBase(value, base)
	if !ok {
		return 0, numSyntax("ParseInt", original)
	}
	if digits == "" {
		return 0, numSyntax("ParseInt", original)
	}

	limit := maxSignedMagnitude(bits, negative)
	magnitude, rangeErr := parseUnsignedDigits(digits, parseBase, limit)
	if rangeErr != nil {
		if negative {
			if bits == 64 {
				return -1 << 63, numRange("ParseInt", original)
			}
			return -int64(uint64(1) << (bits - 1)), numRange("ParseInt", original)
		}

		if bits == 64 {
			return int64(^uint64(0) >> 1), numRange("ParseInt", original)
		}
		return int64((uint64(1) << (bits - 1)) - 1), numRange("ParseInt", original)
	}

	if negative {
		if bits == 64 && magnitude == uint64(1)<<63 {
			return -1 << 63, nil
		}
		return -int64(magnitude), nil
	}

	return int64(magnitude), nil
}

func ParseUint(value string, base int, bitSize int) (uint64, error) {
	original := value
	if value == "" {
		return 0, numSyntax("ParseUint", original)
	}
	if value[0] == '+' {
		value = value[1:]
	}
	if value == "" || value[0] == '-' {
		return 0, numSyntax("ParseUint", original)
	}

	bits := normalizeBitSize(bitSize)
	parseBase, digits, ok := normalizeBase(value, base)
	if !ok {
		return 0, numSyntax("ParseUint", original)
	}
	if digits == "" {
		return 0, numSyntax("ParseUint", original)
	}

	limit := maxUnsigned(bits)
	parsed, err := parseUnsignedDigits(digits, parseBase, limit)
	if err != nil {
		return limit, numRange("ParseUint", original)
	}

	return parsed, nil
}

func appendUnsigned(dst []byte, value uint64, base uint) []byte {
	if value == 0 {
		return append(dst, '0')
	}

	var scratch [64]byte
	index := len(scratch)
	for value > 0 {
		quotient, remainder := divModUint64(value, uint32(base))
		index--
		scratch[index] = lowerDigits[remainder]
		value = quotient
	}

	return append(dst, scratch[index:]...)
}

func normalizeFormatBase(base int) uint {
	if base < 2 || base > 36 {
		return 10
	}

	return uint(base)
}

func normalizeBase(value string, base int) (normalized uint, digits string, ok bool) {
	switch {
	case base == 0:
		switch {
		case hasPrefixFold(value, "0x"):
			return 16, value[2:], true
		case hasPrefixFold(value, "0b"):
			return 2, value[2:], true
		case hasPrefixFold(value, "0o"):
			return 8, value[2:], true
		case len(value) > 1 && value[0] == '0':
			return 8, value[1:], true
		default:
			return 10, value, true
		}
	case base == 16 && hasPrefixFold(value, "0x"):
		return 16, value[2:], true
	case base == 2 && hasPrefixFold(value, "0b"):
		return 2, value[2:], true
	case base == 8 && hasPrefixFold(value, "0o"):
		return 8, value[2:], true
	case base >= 2 && base <= 36:
		return uint(base), value, true
	default:
		return 0, "", false
	}
}

func parseUnsignedDigits(value string, base uint, limit uint64) (uint64, error) {
	current := uint64(0)
	base64 := uint64(base)
	limitQuotient, limitRemainder := divModUint64(limit, uint32(base))

	for index := 0; index < len(value); index++ {
		digit, ok := digitValue(value[index])
		if !ok || digit >= base64 {
			return 0, ErrSyntax
		}
		if current > limitQuotient || (current == limitQuotient && digit > uint64(limitRemainder)) {
			return limit, ErrRange
		}

		current = current*base64 + digit
	}

	return current, nil
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

func maxUnsigned(bits uint) uint64 {
	if bits >= 64 {
		return ^uint64(0)
	}

	return (uint64(1) << bits) - 1
}

func normalizeBitSize(bitSize int) uint {
	switch {
	case bitSize == 0:
		return IntSize
	case bitSize < 0:
		return IntSize
	case bitSize > 64:
		return 64
	default:
		return uint(bitSize)
	}
}

func numSyntax(function string, value string) error {
	return &NumError{
		Func: function,
		Num:  value,
		Err:  ErrSyntax,
	}
}

func numRange(function string, value string) error {
	return &NumError{
		Func: function,
		Num:  value,
		Err:  ErrRange,
	}
}

func hasPrefixFold(value string, prefix string) bool {
	if len(value) < len(prefix) {
		return false
	}

	return equalFoldASCII(value[:len(prefix)], prefix)
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

func digitValue(value byte) (uint64, bool) {
	switch {
	case value >= '0' && value <= '9':
		return uint64(value - '0'), true
	case value >= 'a' && value <= 'z':
		return uint64(value-'a') + 10, true
	case value >= 'A' && value <= 'Z':
		return uint64(value-'A') + 10, true
	}

	return 0, false
}

func quote(value string) string {
	return `"` + value + `"`
}

func divModUint64(value uint64, divisor uint32) (uint64, uint8) {
	quotient := uint64(0)
	remainder := uint64(0)
	divisor64 := uint64(divisor)

	for shift := uint(64); shift > 0; shift-- {
		remainder = (remainder << 1) | ((value >> (shift - 1)) & 1)
		if remainder >= divisor64 {
			remainder -= divisor64
			quotient |= uint64(1) << (shift - 1)
		}
	}

	return quotient, uint8(remainder)
}

var lowerDigits = [...]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
}
