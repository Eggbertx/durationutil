package durationutil

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

var (
	ErrEmptyDurationString   = errors.New("empty duration string")
	ErrInvalidDurationString = errors.New("invalid duration string")
	durationRegexp           = regexp.MustCompile(`^((\d+)\s?ye?a?r?s?)?\s?((\d+)\s?mon?t?h?s?)?\s?((\d+)\s?we?e?k?s?)?\s?((\d+)\s?da?y?s?)?\s?((\d+)\s?ho?u?r?s?)?\s?((\d+)\s?mi?n?u?t?e?s?)?\s?((\d+)\s?s?e?c?o?n?d?s?)?$`)
)

// ParseLongerDuration parses the given string into a duration and returns any errors.
// Based on TinyBoard's parse_time function
func ParseLongerDuration(str string) (time.Duration, error) {
	if str == "" {
		return 0, ErrEmptyDurationString
	}

	matches := durationRegexp.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		return 0, ErrInvalidDurationString
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
