package durationutil

import (
	"encoding/json"
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
	var expectedErr *InvalidDurationStringError
	if duration, err = ParseLongerDuration(durStr); !errors.As(err, &expectedErr) || expectedErr.Value != durStr {
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

type testStruct struct {
	Years  []ExtendedDuration `json:",omitempty"`
	Months []ExtendedDuration `json:",omitempty"`
	Weeks  []ExtendedDuration `json:",omitempty"`
	Zero   ExtendedDuration
}

func TestEncodingDecodingJSON(t *testing.T) {
	var tc testStruct
	tcExpect := testStruct{
		Years:  []ExtendedDuration{Year, 2 * Year, 3 * Year, 4 * Year, 5*Year + 11*Month + 29*Day + 23*ExtendedDuration(time.Hour) + 59*ExtendedDuration(time.Minute) + 59*ExtendedDuration(time.Second)},
		Months: []ExtendedDuration{Month, Month, Month, Month},
		Weeks:  []ExtendedDuration{Week, Week, Week, Week + 2*ExtendedDuration(time.Minute)},
	}
	err := json.Unmarshal([]byte(`{
		"Years": ["1y", "2 y", "3 year", "4year", "5 years 11 months 29 days 23 hours 59 minutes 59 seconds"],
		"Months": ["1mo", "1 mo", "1 month", "1month"],
		"Weeks": ["1w", "1 w", "1 week", "1week 2 minutes"]
	}`), &tc)

	t.Run("Normal JSON test", func(t *testing.T) {
		if err != nil {
			t.Fatalf("Error decoding JSON: %v", err)
		}
		if len(tc.Years) != len(tcExpect.Years) {
			t.Fatalf("Expected %d years, got %d", len(tcExpect.Years), len(tc.Years))
		}
		for i := range tc.Years {
			if tc.Years[i] != tcExpect.Years[i] {
				t.Fatalf("Expected year %d to be %v, got %v", i, tcExpect.Years[i], tc.Years[i])
			}
		}
		if len(tc.Months) != len(tcExpect.Months) {
			t.Fatalf("Expected %d months, got %d", len(tcExpect.Months), len(tc.Months))
		}
		for i := range tc.Months {
			if tc.Months[i] != tcExpect.Months[i] {
				t.Fatalf("Expected month %d to be %v, got %v", i, tcExpect.Months[i], tc.Months[i])
			}
		}
		if len(tc.Weeks) != len(tcExpect.Weeks) {
			t.Fatalf("Expected %d weeks, got %d", len(tcExpect.Weeks), len(tc.Weeks))
		}
		for i := range tc.Weeks {
			if tc.Weeks[i] != tcExpect.Weeks[i] {
				t.Fatalf("Expected week %d to be %v, got %v", i, tcExpect.Weeks[i], tc.Weeks[i])
			}
		}
	})

	t.Run("JSON encoding test", func(t *testing.T) {
		tcEncoded, err := json.Marshal(tc)
		if err != nil {
			t.Fatalf("Error encoding JSON: %v", err)
		}
		tcEncodedExpect := []byte(`{"Years":["1y","2y","3y","4y","5y51w2d23h59m59s"],"Months":["4w2d","4w2d","4w2d","4w2d"],"Weeks":["1w","1w","1w","1w2m0s"],"Zero":""}`)
		if tcEncoded == nil || string(tcEncoded) != string(tcEncodedExpect) {
			t.Fatalf("Expected encoded JSON to be %s, got %s", string(tcEncodedExpect), string(tcEncoded))
		}
	})

	t.Run("Invalid types in JSON", func(t *testing.T) {
		var arr []ExtendedDuration
		err = json.Unmarshal([]byte(`[4.2, "1 year", "2 months"]`), &arr)
		expect := &InvalidDurationStringError{"4.2"}
		if !errors.Is(err, expect) {
			t.Fatalf("Expected error to be of type %v, got %v", err, expect)
		}
	})

	t.Run("Handle nil ExtendedDuration pointer", func(t *testing.T) {
		var ed *ExtendedDuration
		jsonData, err := json.Marshal(ed)
		if err != nil {
			t.Fatalf("Error marshaling nil ExtendedDuration: %v", err)
		}
		if string(jsonData) != "null" {
			t.Fatalf("Expected marshaled nil ExtendedDuration to be null, got %s", string(jsonData))
		}

		var edUnmarshal ExtendedDuration
		err = json.Unmarshal([]byte("null"), &edUnmarshal)
		if !errors.Is(err, &InvalidDurationStringError{Value: "null"}) {
			t.Fatalf("Expected %v when unmarshaling null, got %v", &InvalidDurationStringError{Value: "null"}, err)
		}
	})

	t.Run("empty string handling in JSON", func(t *testing.T) {
		var tc testStruct
		err = json.Unmarshal([]byte(`{"Zero":""}`), &tc)
		if err != nil {
			t.Fatalf("Error unmarshaling empty string: %v", err)
		}
		if tc.Zero != 0 {
			t.Fatalf("Expected unmarshaled empty string to be 0, got %v", tc.Zero)
		}

		jsonData, err := json.Marshal(tc)
		if err != nil {
			t.Fatalf("Error marshaling zero ExtendedDuration: %v", err)
		}
		if string(jsonData) != `{"Zero":""}` {
			t.Fatalf("Expected marshaled zero ExtendedDuration to be empty string, got %s", string(jsonData))
		}
	})
}
