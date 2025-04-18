package main

import (
	"testing"
)

func TestTimeParse(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool // true if we expect successful parsing
	}{
		{"2023-11-15 10:30:45", true},
		{"15 Nov 23 14:30 EST", true},
		{"2023-11-15T10:30:45Z", true},
		{"2023-11-15 10:30:45 UTC", true},
		{"Wed, 15 Nov 2023 10:30:45 GMT", true}, // RFC1123 format, not in your list
		{"15 Nov 2023 10:30:45 +0000", true},    // Another common format not in your list
		{"invalid date format", false},
	}

	for i, tc := range testCases {
		result := time_parse(tc.input)
		isZero := result.IsZero()

		if tc.expected && isZero {
			t.Errorf("Test case %d: Expected to parse '%s' successfully, but got zero time", i, tc.input)
		} else if !tc.expected && !isZero {
			t.Errorf("Test case %d: Expected to fail parsing '%s', but got %v", i, tc.input, result)
		}

		if !isZero {
			t.Logf("Successfully parsed '%s' to: %v", tc.input, result)
		}
	}
}
