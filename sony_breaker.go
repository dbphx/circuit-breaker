package ecocircuitbreaker

import "github.com/sony/gobreaker/v2"

type SonyCircuitBreaker[T any] struct {
	cb *gobreaker.CircuitBreaker[T]
}

func NewSonyBreaker[T any](cfg *BreakerOptions) IBreaker[T] {
	settings := gobreaker.Settings{
		Name: "sonyCB",
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= uint32(cfg.FailureThreshold)
		},
		Timeout: cfg.OpenTimeout,
		OnStateChange: func(name string, from, to gobreaker.State) {
			if cfg.OnStateChange != nil {
				cfg.OnStateChange(toState(from), toState(to))
			}
		},
	}
	return &SonyCircuitBreaker[T]{
		cb: gobreaker.NewCircuitBreaker[T](settings),
	}
}

func (s *SonyCircuitBreaker[T]) Execute(fn func() (T, error)) (T, error) {
	return s.cb.Execute(fn)
}

func toState(s gobreaker.State) State {
	switch s {
	case gobreaker.StateClosed:
		return Closed
	case gobreaker.StateHalfOpen:
		return HalfOpen
	default:
		return Open
	}
}

func (s *SonyCircuitBreaker[T]) CurrentState() State {
	return toState(s.cb.State())
}
