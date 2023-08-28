package durationutil

import (
	"errors"
	"testing"
	"time"
)

const (
	errorParsingFmt = "Error parsing %q: %s"
	expectedFmt     = "Expected %v, got %v"

	expectedDuration time.Duration = time.Hour*66579 + time.Minute*2 + time.Second*10
)

func TestParseLongerDuration(t *testing.T) {
	durStr := ""
	// test invalid string parsing returns correct errors
	duration, err := ParseLongerDuration(durStr)
	if !errors.Is(err, ErrEmptyDurationString) {
		t.Fatalf("Expected empty duration string to return empty duration string error, got %v", err)
	}

	durStr = "seven years, six months, five days, 4w3d"
	duration, err = ParseLongerDuration(durStr)
	if duration, err = ParseLongerDuration(durStr); !errors.Is(err, ErrInvalidDurationString) {
		t.Fatalf("Expected invalid duration string error, got %v", err)
	}
	if duration != 0 {
		t.Fatalf(expectedFmt, expectedDuration, duration)
	}

	// test valid string parsing
	durStr = "7y6mo5w4d3h2m10s"
	if duration, err = ParseLongerDuration(durStr); err != nil {
		t.Fatalf(errorParsingFmt, durStr, err.Error())
	}
	if duration != expectedDuration {
		t.Fatalf(expectedFmt, expectedDuration, duration)
	}

	durStr = "7year6month5weeks4days3hours2minutes10second"
	if duration, err = ParseLongerDuration(durStr); err != nil {
		t.Fatalf(errorParsingFmt, durStr, err.Error())
	}
	if duration != expectedDuration {
		t.Fatalf(expectedFmt, expectedDuration, duration)
	}

	durStr = "7yrs 6 months 5 weeks 4 days 3 hours 2 minutes 10 seconds"
	if duration, err = ParseLongerDuration(durStr); err != nil {
		t.Fatalf(errorParsingFmt, durStr, err.Error())
	}
	if duration != expectedDuration {
		t.Fatalf(expectedFmt, expectedDuration, duration)
	}

	// test valid duration string with skipped units
	durStr = "1 year 2 days 3 seconds"
	if duration, err = ParseLongerDuration(durStr); err != nil {
		t.Fatalf(errorParsingFmt, durStr, err.Error())
	}
	if duration != time.Hour*24*365+time.Hour*24*2+time.Second*3 {
		t.Fatalf(expectedFmt, time.Hour*24*365+time.Hour*24*2+time.Second*3, duration)
	}
}
