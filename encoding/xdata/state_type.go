package xdata

// StateType indicates whether a check is in the SOFT or HARD state.
type StateType int

const (
	Soft StateType = iota
	Hard
)

// ParseStateType converts a byte string into a StateType.
//
// When an unknown type is passed in, a sensible default value is returned, along with ErrUnknownValue.
func ParseStateType(t []byte) (StateType, error) {
	switch string(t) {
	case "SOFT", "0":
		return Soft, nil
	case "HARD", "1":
		return Hard, nil
	default:
		return Soft, ErrUnknownValue
	}
}

// String returns a string representation of StateType.
//
// An empty string is returned for unknown types.
func (t StateType) String() string {
	switch t {
	case Soft:
		return "SOFT"
	case Hard:
		return "HARD"
	default:
		return ""
	}
}
