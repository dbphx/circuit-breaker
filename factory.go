package ecocircuitbreaker

import "time"

func NewBreaker[T any](opts ...BreakerOption) IBreaker[T] {
	cfg := &BreakerOptions{
		Type:                 EcoBreakerType,
		FailureThreshold:     3,
		SuccessThreshold:     1,
		OpenTimeout:          5 * time.Second,
		FailureResetDuration: 2 * time.Minute,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	switch cfg.Type {
	case SonyBreakerType:
		return NewSonyBreaker[T](cfg)
	default:
		return NewEcoBreaker[T](cfg)
	}
}
