package xdata

import (
	"strings"
)

// CheckType identifies whether a check is actively checked by Nagios (Active),
// or recieves data from an external system (Passive).
type CheckType int

const (
	Active CheckType = iota
	Passive
)

// ParseCheckType converts a byte string into a CheckType.
//
// When an unknown type is passed in, a sensible default is returned along with
// ErrUnknownAcknowledgementType.
func ParseCheckType(t []byte) (CheckType, error) {
	switch strings.TrimSpace(strings.ToLower(string(t))) {
	case "active", "0":
		return Active, nil
	case "passive", "1":
		return Passive, nil
	default:
		return Active, ErrUnknownValue
	}
}

// String returns a string representation of the CheckType.
//
// An empty string is returned for unknown types.
func (t CheckType) String() string {
	switch t {
	case Active:
		return "Active"
	case Passive:
		return "Passive"
	default:
		return ""
	}
}
