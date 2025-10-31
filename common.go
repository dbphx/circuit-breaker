package ecocircuitbreaker

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

func (s State) String() string {
	switch s {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case HalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

func (s State) IsClosed() bool {
	return s == Closed
}

func (s State) IsOpen() bool {
	return s == Open
}

func (s State) IsHalfOpen() bool {
	return s == HalfOpen
}
