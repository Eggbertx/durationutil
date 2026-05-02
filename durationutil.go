package durationutil

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Day   ExtendedDuration = ExtendedDuration(24 * time.Hour)
	Week                   = 7 * Day
	Month                  = 30 * Day
	Year                   = 365 * Day
)

var (
	ErrEmptyDurationString = errors.New("empty duration string")

	// ErrInvalidDurationString represents an invalid duration string
	//
	// Deprecated: ParseLongerDuration now returns an InvalidDurationStringError typed error with the specific invalid value
	ErrInvalidDurationString error = &InvalidDurationStringError{""}
	durationRegexp                 = regexp.MustCompile(`^((\d+)\s?ye?a?r?s?)?\s?((\d+)\s?mon?t?h?s?)?\s?((\d+)\s?we?e?k?s?)?\s?((\d+)\s?da?y?s?)?\s?((\d+)\s?ho?u?r?s?)?\s?((\d+)\s?mi?n?u?t?e?s?)?\s?((\d+)\s?s?e?c?o?n?d?s?)?$`)
)

// ExtendedDuration is a wrapper around time.Duration that can be encoded into and decoded from JSON using durations
// that include years, months, weeks, and days. If the string is blank, the duration will be set to 0, instead of
// returning an error like ParseLongerDuration
type ExtendedDuration time.Duration

func (ed *ExtendedDuration) UnmarshalJSON(ba []byte) (err error) {
	baStr := strings.Trim(string(ba), `"`)
	if baStr == "" {
		*ed = 0
		return nil
	}
	dur, err := ParseLongerDuration(baStr)
	if err == nil {
		*ed = ExtendedDuration(dur)
	}
	return err
}

func (ed ExtendedDuration) String() string {
	var years int
	var weeks int
	var days int

	trimmed := ed
	for trimmed >= Year {
		years++
		trimmed -= Year
	}
	for trimmed >= Week {
		weeks++
		trimmed -= Week
	}
	for trimmed >= Day {
		days++
		trimmed -= Day
	}
	var strOut string
	if years > 0 {
		strOut += strconv.Itoa(years) + "y"
	}
	if weeks > 0 {
		strOut += strconv.Itoa(weeks) + "w"
	}
	if days > 0 {
		strOut += strconv.Itoa(days) + "d"
	}
	if trimmed > 0 {
		strOut += time.Duration(trimmed).String()
	}
	return strOut
}

func (ed ExtendedDuration) MarshalJSON() ([]byte, error) {
	if ed == 0 {
		return []byte(`""`), nil
	}

	return []byte(`"` + ed.String() + `"`), nil
}

type InvalidDurationStringError struct {
	Value string
}

func (e *InvalidDurationStringError) Error() string {
	if e == nil || e.Value == "" {
		return "invalid duration string"
	}
	return "invalid duration string: " + e.Value
}

func (e *InvalidDurationStringError) Is(err error) bool {
	if err == nil {
		return false
	}
	asE, ok := err.(*InvalidDurationStringError)
	return ok && asE.Value == e.Value
}

// ParseLongerDuration parses the given string into a duration and returns any errors.
// Based on TinyBoard's parse_time function
func ParseLongerDuration(str string) (time.Duration, error) {
	if str == "" {
		return 0, ErrEmptyDurationString
	}

	matches := durationRegexp.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return 0, &InvalidDurationStringError{Value: str}
	}

	var duration int
	if matches[0][2] != "" {
		years, _ := strconv.Atoi(matches[0][2])
		duration += years * 60 * 60 * 24 * 365
	}
	if matches[0][4] != "" {
		months, _ := strconv.Atoi(matches[0][4])
		duration += months * 60 * 60 * 24 * 30
	}
	if matches[0][6] != "" {
		weeks, _ := strconv.Atoi(matches[0][6])
		duration += weeks * 60 * 60 * 24 * 7
	}
	if matches[0][8] != "" {
		days, _ := strconv.Atoi(matches[0][8])
		duration += days * 60 * 60 * 24
	}
	if matches[0][10] != "" {
		hours, _ := strconv.Atoi(matches[0][10])
		duration += hours * 60 * 60
	}
	if matches[0][12] != "" {
		minutes, _ := strconv.Atoi(matches[0][12])
		duration += minutes * 60
	}
	if matches[0][14] != "" {
		seconds, _ := strconv.Atoi(matches[0][14])
		duration += seconds
	}
	return time.Duration(duration) * time.Second, nil
}
