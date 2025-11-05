package ecocircuitbreaker

import "time"

type BreakerType int

const (
	EcoBreakerType BreakerType = iota
	SonyBreakerType
)

type BreakerOptions struct {
	Type                 BreakerType
	FailureThreshold     int
	SuccessThreshold     int
	FailureResetDuration time.Duration
	OpenTimeout          time.Duration
	OnStateChange        func(old, new State)
}

type BreakerOption func(*BreakerOptions)

func WithType(t BreakerType) BreakerOption {
	return func(cfg *BreakerOptions) { cfg.Type = t }
}

func WithFailureThreshold(n int) BreakerOption {
	return func(cfg *BreakerOptions) {
		if n > 0 {
			cfg.FailureThreshold = n
		}
	}
}

func WithSuccessThreshold(n int) BreakerOption {
	return func(cfg *BreakerOptions) {
		if n > 0 {
			cfg.SuccessThreshold = n
		}
	}
}

func WithOpenTimeout(d time.Duration) BreakerOption {
	return func(cfg *BreakerOptions) {
		if d > 0 {
			cfg.OpenTimeout = d
		}
	}
}

func WithOnStateChange(fn func(old, new State)) BreakerOption {
	return func(cfg *BreakerOptions) {
		cfg.OnStateChange = fn
	}
}

func WithFailureResetDuration(d time.Duration) BreakerOption {
	return func(cfg *BreakerOptions) {
		if d > 0 {
			cfg.FailureResetDuration = d
		}
	}
}
