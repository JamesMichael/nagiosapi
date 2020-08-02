package xdata

import "strings"

// ServiceState represents the current statusdata of a service check.
type (
	ServiceState int
)

const (
	Ok ServiceState = iota
	Warning
	Critical
	Unknown
)

// ParseServiceState converts a byte string into a ServiceState.
//
// When an unknown type is passed in, a sensible default value is returned,
// along with ErrUnknownValue.
func ParseServiceState(s []byte) (ServiceState, error) {
	switch strings.TrimSpace(strings.ToLower(string(s))) {
	case "ok", "0":
		return Ok, nil
	case "warning", "1":
		return Warning, nil
	case "critical", "2":
		return Critical, nil
	case "unknown", "3":
		return Unknown, nil
	default:
		return Ok, ErrUnknownValue
	}
}

// String returns a string representation of ServiceState.
//
// An empty string is returned for unknown types.
func (s ServiceState) String() string {
	switch s {
	case Ok:
		return "OK"
	case Warning:
		return "WARNING"
	case Critical:
		return "CRITICAL"
	case Unknown:
		return "UNKNOWN"
	default:
		return ""
	}
}
