package time

import "kos"

type Duration int64

const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

type Month int

const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

type Time struct {
	unixSeconds  int64
	nanosecond   int32
	monotonicNS  int64
	hasMonotonic bool
}

const (
	nanosecondsPerSecond      = int64(1000000000)
	nanosecondsPerCentisecond = int64(10000000)
	secondsPerMinute          = int64(60)
	secondsPerHour            = int64(60 * 60)
	secondsPerDay             = int64(24 * 60 * 60)
	daysPer400Years           = int64(146097)
	unixToCivilEpochDays      = int64(719468)
	maxInt64                  = int64(1<<63 - 1)
	minInt64                  = int64(-1 << 63)
	maxDurationSeconds        = maxInt64 / nanosecondsPerSecond
	minDurationSeconds        = minInt64 / nanosecondsPerSecond
)

func Now() Time {
	startClock := kos.SystemTime()
	date := kos.SystemDate()
	endClock := kos.SystemTime()
	if clockSeconds(endClock) < clockSeconds(startClock) {
		date = kos.SystemDate()
	}

	seconds := unixFromCivil(
		int64(expandClockYear(date.Year)),
		int64(date.Month),
		int64(date.Day),
		int64(endClock.Hour),
		int64(endClock.Minute),
		int64(endClock.Second),
	)

	return Time{
		unixSeconds:  seconds,
		nanosecond:   0,
		monotonicNS:  int64(kos.UptimeNanoseconds()),
		hasMonotonic: true,
	}
}

func Unix(sec int64, nsec int64) Time {
	sec, nsec = normalizeUnix(sec, nsec)
	return Time{
		unixSeconds: sec,
		nanosecond:  int32(nsec),
	}
}

func Sleep(duration Duration) {
	if duration <= 0 {
		return
	}

	centiseconds, remainder := divModUint64(unsignedAbsInt64(int64(duration)), uint32(nanosecondsPerCentisecond))
	if remainder != 0 {
		centiseconds++
	}

	for centiseconds > 0 {
		chunk := centiseconds
		if chunk > uint64(^uint32(0)) {
			chunk = uint64(^uint32(0))
		}

		kos.SleepCentiseconds(uint32(chunk))
		centiseconds -= chunk
	}
}

func Since(value Time) Duration {
	return Now().Sub(value)
}

func (value Time) Add(duration Duration) Time {
	result := Unix(value.unixSeconds, int64(value.nanosecond)+int64(duration))
	if value.hasMonotonic {
		result.monotonicNS = value.monotonicNS + int64(duration)
		result.hasMonotonic = true
	}

	return result
}

func (value Time) Sub(other Time) Duration {
	if value.hasMonotonic && other.hasMonotonic {
		return clampDurationParts(0, value.monotonicNS-other.monotonicNS)
	}

	return clampDurationParts(value.unixSeconds-other.unixSeconds, int64(value.nanosecond)-int64(other.nanosecond))
}

func (value Time) Before(other Time) bool {
	return value.compare(other) < 0
}

func (value Time) After(other Time) bool {
	return value.compare(other) > 0
}

func (value Time) Equal(other Time) bool {
	return value.compare(other) == 0
}

func (value Time) IsZero() bool {
	return value.unixSeconds == 0 && value.nanosecond == 0
}

func (value Time) Unix() int64 {
	return value.unixSeconds
}

func (value Time) Nanosecond() int {
	return int(value.nanosecond)
}

func (value Time) Second() int {
	_, _, _, _, _, second := value.dateTime()
	return second
}

func (value Time) Minute() int {
	_, _, _, _, minute, _ := value.dateTime()
	return minute
}

func (value Time) Hour() int {
	_, _, _, hour, _, _ := value.dateTime()
	return hour
}

func (value Time) Day() int {
	_, _, day, _, _, _ := value.dateTime()
	return day
}

func (value Time) Month() Month {
	_, month, _, _, _, _ := value.dateTime()
	return month
}

func (value Time) Year() int {
	year, _, _, _, _, _ := value.dateTime()
	return year
}

func (value Time) compare(other Time) int {
	if value.hasMonotonic && other.hasMonotonic {
		switch {
		case value.monotonicNS < other.monotonicNS:
			return -1
		case value.monotonicNS > other.monotonicNS:
			return 1
		default:
			return 0
		}
	}

	switch {
	case value.unixSeconds < other.unixSeconds:
		return -1
	case value.unixSeconds > other.unixSeconds:
		return 1
	case value.nanosecond < other.nanosecond:
		return -1
	case value.nanosecond > other.nanosecond:
		return 1
	default:
		return 0
	}
}

func (value Time) dateTime() (year int, month Month, day int, hour int, minute int, second int) {
	days, daySeconds := divModFloorInt64(value.unixSeconds, uint32(secondsPerDay))

	year, month, day = civilFromDays(days)
	hourQuotient, hourRemainder := divModUint64(uint64(daySeconds), uint32(secondsPerHour))
	minuteQuotient, secondRemainder := divModUint64(uint64(hourRemainder), uint32(secondsPerMinute))
	hour = int(hourQuotient)
	minute = int(minuteQuotient)
	second = int(secondRemainder)
	return
}

func normalizeUnix(seconds int64, nanoseconds int64) (int64, int64) {
	if nanoseconds >= 0 {
		quotient, remainder := divModUint64(uint64(nanoseconds), uint32(nanosecondsPerSecond))
		seconds += int64(quotient)
		nanoseconds = int64(remainder)
	} else {
		quotient, remainder := divModUint64(unsignedAbsInt64(nanoseconds), uint32(nanosecondsPerSecond))
		seconds -= int64(quotient)
		nanoseconds = -int64(remainder)
	}
	if nanoseconds < 0 {
		nanoseconds += nanosecondsPerSecond
		seconds--
	}

	return seconds, nanoseconds
}

func clampDurationParts(seconds int64, nanoseconds int64) Duration {
	if seconds > maxDurationSeconds {
		return Duration(maxInt64)
	}
	if seconds < minDurationSeconds {
		return Duration(minInt64)
	}

	total := seconds * nanosecondsPerSecond
	if nanoseconds > 0 && total > maxInt64-nanoseconds {
		return Duration(maxInt64)
	}
	if nanoseconds < 0 && total < minInt64-nanoseconds {
		return Duration(minInt64)
	}

	return Duration(total + nanoseconds)
}

func clockSeconds(value kos.ClockTime) int64 {
	return int64(value.Hour)*secondsPerHour +
		int64(value.Minute)*secondsPerMinute +
		int64(value.Second)
}

func expandClockYear(year byte) int {
	return 2000 + int(year)
}

func unixFromCivil(year int64, month int64, day int64, hour int64, minute int64, second int64) int64 {
	days := daysFromCivil(year, month, day)
	return days*secondsPerDay + hour*secondsPerHour + minute*secondsPerMinute + second
}

func daysFromCivil(year int64, month int64, day int64) int64 {
	if month <= 2 {
		year--
	}

	era, _ := divModFloorInt64(year, 400)
	yearOfEra := uint32(year - era*400)
	monthPrime := uint32(month)
	if monthPrime > 2 {
		monthPrime -= 3
	} else {
		monthPrime += 9
	}

	dayOfYear := ((153 * monthPrime) + 2) / 5
	dayOfYear += uint32(day) - 1
	dayOfEra := uint64(yearOfEra)*365 + uint64(yearOfEra/4) - uint64(yearOfEra/100) + uint64(dayOfYear)
	return era*daysPer400Years + int64(dayOfEra) - unixToCivilEpochDays
}

func civilFromDays(days int64) (year int, month Month, day int) {
	days += unixToCivilEpochDays
	era, _ := divModFloorInt64(days, uint32(daysPer400Years))
	dayOfEra := uint32(days - era*daysPer400Years)
	yearOfEra := (dayOfEra - dayOfEra/1460 + dayOfEra/36524 - dayOfEra/146096) / 365
	yearValue := int64(yearOfEra) + era*400
	dayOfYear := dayOfEra - (365*yearOfEra + yearOfEra/4 - yearOfEra/100)
	monthPrime := (5*dayOfYear + 2) / 153

	day = int(dayOfYear - ((153*monthPrime+2)/5) + 1)
	if monthPrime < 10 {
		month = Month(monthPrime + 3)
	} else {
		month = Month(monthPrime - 9)
		yearValue++
	}
	year = int(yearValue)
	return
}

func unsignedAbsInt64(value int64) uint64 {
	if value >= 0 {
		return uint64(value)
	}

	return uint64(^value) + 1
}

func divModUint64(value uint64, divisor uint32) (uint64, uint32) {
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

	return quotient, uint32(remainder)
}

func divModFloorInt64(value int64, divisor uint32) (int64, uint32) {
	if value >= 0 {
		quotient, remainder := divModUint64(uint64(value), divisor)
		return int64(quotient), remainder
	}

	quotient, remainder := divModUint64(unsignedAbsInt64(value), divisor)
	if remainder == 0 {
		return -int64(quotient), 0
	}

	return -int64(quotient) - 1, divisor - remainder
}
