package xdata

import (
	"errors"
	"testing"
)

func TestParseServiceState(t *testing.T) {
	tests := []struct {
		input    string
		expected ServiceState
	}{
		{"OK", Ok},
		{"WARNING", Warning},
		{"CRITICAL", Critical},
		{"UNKNOWN", Unknown},
		{"0", Ok},
		{"1", Warning},
		{"2", Critical},
		{"3", Unknown},
	}

	for _, test := range tests {
		st, err := ParseServiceState([]byte(test.input))
		if err != nil {
			t.Errorf("unable to parse service state '%s': %s", test.input, err)
			continue
		}

		if st != test.expected {
			t.Errorf("parse returned incorrect output, got: '%d', want: '%d", st, test.expected)
		}
	}

	st, err := ParseServiceState([]byte("invalid"))
	if err == nil {
		t.Errorf("expected error when passing in unknown type")
	} else if !errors.Is(err, ErrUnknownValue) {
		t.Errorf("got unknown error when passing in unknown type: %s", err)
	}

	if st != None {
		t.Errorf("parse returned incorrect output, got: '%d', want: '%d'", st, None)
	}
}

func TestServiceState_String(t *testing.T) {
	tests := []struct {
		input    ServiceState
		expected string
	}{
		{Ok, "OK"},
		{Warning, "WARNING"},
		{Critical, "CRITICAL"},
		{Unknown, "UNKNOWN"},
		{ServiceState(100), ""},
	}

	for _, test := range tests {
		s := test.input.String()
		if s != test.expected {
			t.Errorf("unexpected string, got: '%s', want: '%s'", s, test.expected)
		}
	}
}
