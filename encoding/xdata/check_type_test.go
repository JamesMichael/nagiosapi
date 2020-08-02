package xdata

import (
	"errors"
	"testing"
)

func TestParseCheckType(t *testing.T) {
	tests := []struct {
		input    string
		expected CheckType
	}{
		{"Active", Active},
		{"Passive", Passive},
		{"0", Active},
		{"1", Passive},
	}

	for _, test := range tests {
		ct, err := ParseCheckType([]byte(test.input))
		if err != nil {
			t.Errorf("unable to parse check type '%s': %s", test.input, err)
			continue
		}

		if ct != test.expected {
			t.Errorf("parse returned incorrect output, got: '%d', want: '%d", ct, test.expected)
		}
	}

	ct, err := ParseCheckType([]byte("unknown"))
	if err == nil {
		t.Errorf("expected error when passing in unknown type")
	} else if !errors.Is(err, ErrUnknownValue) {
		t.Errorf("got unknown error when passing in unknown type: %s", err)
	}

	if ct != None {
		t.Errorf("parse returned incorrect output, got: '%d', want: '%d'", ct, None)
	}
}

func TestCheckType_String(t *testing.T) {
	tests := []struct {
		input    CheckType
		expected string
	}{
		{Active, "Active"},
		{Passive, "Passive"},
		{CheckType(100), ""},
	}

	for _, test := range tests {
		s := test.input.String()
		if s != test.expected {
			t.Errorf("unexpected string, got: '%s', want: '%s'", s, test.expected)
		}
	}

}
