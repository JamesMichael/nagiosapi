package xdata

import (
	"strings"
)

// AcknowledgementType identifies whether the host/service check has been acknowledged and if so, whether the
// acknowledgement will persist (Sticky) or not (Normal) once the check turns Ok.
type AcknowledgementType int

const (
	None = iota
	Normal
	Sticky
)

// ParseAcknowledgementType converts a byte string into an AcknowledgementType.
//
// When an unknown type is passed in, a sensible default is returned along with ErrUnknownValue.
func ParseAcknowledgementType(t []byte) (AcknowledgementType, error) {
	switch strings.TrimSpace(strings.ToLower(string(t))) {
	case "none", "0":
		return None, nil
	case "normal", "1":
		return Normal, nil
	case "sticky", "2":
		return Sticky, nil
	default:
		return None, ErrUnknownValue
	}
}

// String returns a string representation of the AcknowledgementType.
//
// An empty string is returned for unknown types.
func (t AcknowledgementType) String() string {
	switch t {
	case None:
		return "None"
	case Normal:
		return "Normal"
	case Sticky:
		return "Sticky"
	default:
		return ""
	}
}
