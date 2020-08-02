package xdata

import (
	"errors"
	"testing"
)

func TestParseHostState(t *testing.T) {
	tests := []struct {
		input    string
		expected HostState
	}{
		{"UP", Up},
		{"DOWN", Down},
		{"UNREACHABLE", Unreachable},
		{"0", Up},
		{"1", Down},
		{"2", Unreachable},
	}

	for _, test := range tests {
		st, err := ParseHostState([]byte(test.input))
		if err != nil {
			t.Errorf("unable to parse host state '%s': %s", test.input, err)
			continue
		}

		if st != test.expected {
			t.Errorf("parse returned incorrect output, got: '%d', want: '%d", st, test.expected)
		}
	}

	st, err := ParseHostState([]byte("unknown"))
	if err == nil {
		t.Errorf("expected error when passing in unknown type")
	} else if !errors.Is(err, ErrUnknownValue) {
		t.Errorf("got unknown error when passing in unknown type: %s", err)
	}

	if st != None {
		t.Errorf("parse returned incorrect output, got: '%d', want: '%d'", st, None)
	}
}

func TestHostState_String(t *testing.T) {
	tests := []struct {
		input    HostState
		expected string
	}{
		{Up, "UP"},
		{Down, "DOWN"},
		{Unreachable, "UNREACHABLE"},
		{HostState(100), ""},
	}

	for _, test := range tests {
		s := test.input.String()
		if s != test.expected {
			t.Errorf("unexpected string, got: '%s', want: '%s'", s, test.expected)
		}
	}
}
