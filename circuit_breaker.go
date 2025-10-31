package ecocircuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var ErrOpen = errors.New("circuit breaker is open")

type CircuitBreaker[T any] struct {
	mu sync.Mutex

	state      State
	failures   int
	successes  int
	lastOpened time.Time

	failureThreshold int
	successThreshold int
	openTimeout      time.Duration
	onStateChange    func(old, new State)
}

func NewEcoBreaker[T any](cfg *BreakerOptions) IBreaker[T] {
	return &CircuitBreaker[T]{
		state:            Closed,
		failureThreshold: cfg.FailureThreshold,
		successThreshold: cfg.SuccessThreshold,
		openTimeout:      cfg.OpenTimeout,
		onStateChange:    cfg.OnStateChange,
	}
}

func (cb *CircuitBreaker[T]) Execute(fn func() (T, error)) (T, error) {
	cb.mu.Lock()
	state := cb.currentStateNonBlocking()
	if state == Open {
		cb.mu.Unlock()
		var zero T
		return zero, ErrOpen
	}
	cb.mu.Unlock()

	val, err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case Closed:
		if err != nil {
			cb.failures++
			if cb.failures >= cb.failureThreshold {
				cb.setState(Open)
				cb.lastOpened = time.Now()
			}
		} else {
			cb.failures = 0
		}
	case HalfOpen:
		if err != nil {
			cb.setState(Open)
			cb.lastOpened = time.Now()
		} else {
			cb.successes++
			if cb.successes >= cb.successThreshold {
				cb.setState(Closed)
			}
		}
	}
	return val, err
}

func (cb *CircuitBreaker[T]) CurrentState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == Open && time.Since(cb.lastOpened) > cb.openTimeout {
		cb.setState(HalfOpen)
	}
	return cb.state
}

func (cb *CircuitBreaker[T]) currentStateNonBlocking() State {
	if cb.state == Open && time.Since(cb.lastOpened) > cb.openTimeout {
		cb.setState(HalfOpen)
	}

	return cb.state
}

func (cb *CircuitBreaker[T]) setState(new State) {
	old := cb.state
	if old == new {
		return
	}
	cb.state = new
	cb.failures = 0
	cb.successes = 0
	if cb.onStateChange != nil {
		go cb.onStateChange(old, new)
	}
}
