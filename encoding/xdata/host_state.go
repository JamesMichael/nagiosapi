package xdata

import (
	"strings"
)

// HostState represents the current statusdata of a host check.
//
// Depending on the type of host check, '1' can either represent 'DOWN/Unreachable' (raw), or
// just 'DOWN' (processed).
type (
	HostState int
)

const (
	Up HostState = iota
	Down
	Unreachable
)

// ParseHostState converts a byte string into a HostState.
//
// When an unknown type is passed in, a sensible default value is returned, along with
// ErrUnknownHostState.
func ParseHostState(s []byte) (HostState, error) {
	switch strings.TrimSpace(strings.ToLower(string(s))) {
	case "up", "0":
		return Up, nil
	case "down", "1":
		return Down, nil
	case "unreachable", "2":
		return Unreachable, nil
	default:
		return Up, ErrUnknownValue
	}
}

// String returns a string representation of the HostState.
//
// An empty string is returned for unknown types.
func (s HostState) String() string {
	switch s {
	case Up:
		return "UP"
	case Down:
		return "DOWN"
	case Unreachable:
		return "UNREACHABLE"
	default:
		return ""
	}
}
