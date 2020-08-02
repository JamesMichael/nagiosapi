package xdata

import (
	"errors"
	"testing"
)

func TestParseAcknowledgementType(t *testing.T) {
	tests := []struct {
		input    string
		expected AcknowledgementType
	}{
		{"None", None},
		{"Normal", Normal},
		{"Sticky", Sticky},
		{"0", None},
		{"1", Normal},
		{"2", Sticky},
	}

	for _, test := range tests {
		at, err := ParseAcknowledgementType([]byte(test.input))
		if err != nil {
			t.Errorf("unable to parse acknowledgement type '%s': %s", test.input, err)
			continue
		}

		if at != test.expected {
			t.Errorf("parse returned incorrect output, got: '%d', want: '%d", at, test.expected)
		}
	}

	at, err := ParseAcknowledgementType([]byte("unknown"))
	if err == nil {
		t.Errorf("expected error when passing in unknown type")
	} else if !errors.Is(err, ErrUnknownValue) {
		t.Errorf("got unknown error when passing in unknown type: %s", err)
	}

	if at != None {
		t.Errorf("parse returned incorrect output, got: '%d', want: '%d'", at, None)
	}
}

func TestAcknowledgementType_String(t *testing.T) {
	tests := []struct {
		input    AcknowledgementType
		expected string
	}{
		{None, "None"},
		{Normal, "Normal"},
		{Sticky, "Sticky"},
		{AcknowledgementType(100), ""},
	}

	for _, test := range tests {
		s := test.input.String()
		if s != test.expected {
			t.Errorf("unexpected string, got: '%s', want: '%s'", s, test.expected)
		}
	}

}
