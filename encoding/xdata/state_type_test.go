package xdata

import (
	"errors"
	"testing"
)

func TestParseStateType(t *testing.T) {
	tests := []struct {
		input    string
		expected StateType
	}{
		{"SOFT", Soft},
		{"HARD", Hard},
		{"0", Soft},
		{"1", Hard},
	}

	for _, test := range tests {
		st, err := ParseStateType([]byte(test.input))
		if err != nil {
			t.Errorf("unable to parse state type '%s': %s", test.input, err)
			continue
		}

		if st != test.expected {
			t.Errorf("parse returned incorrect output, got: '%d', want: '%d", st, test.expected)
		}
	}

	st, err := ParseStateType([]byte("unknown"))
	if err == nil {
		t.Errorf("expected error when passing in unknown type")
	} else if !errors.Is(err, ErrUnknownValue) {
		t.Errorf("got unknown error when passing in unknown type: %s", err)
	}

	if st != Soft {
		t.Errorf("parse returned incorrect output, got: '%d', want: '%d'", st, None)
	}
}

func TestStateType_String(t *testing.T) {
	tests := []struct {
		input    StateType
		expected string
	}{
		{Soft, "SOFT"},
		{Hard, "HARD"},
		{StateType(100), ""},
	}

	for _, test := range tests {
		s := test.input.String()
		if s != test.expected {
			t.Errorf("unexpected string, got: '%s', want: '%s'", s, test.expected)
		}
	}
}
