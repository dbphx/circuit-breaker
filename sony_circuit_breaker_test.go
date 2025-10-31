package ecocircuitbreaker

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestSonyBreaker_ClosedToOpen(t *testing.T) {
	var mu sync.Mutex
	var transitions []State
	var wg sync.WaitGroup
	wg.Add(1)

	cb := NewSonyBreaker[int](&BreakerOptions{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		OpenTimeout:      100 * time.Millisecond,
		OnStateChange: func(old, new State) {
			mu.Lock()
			defer mu.Unlock()
			wg.Done()
			transitions = append(transitions, new)
		},
	}).(*SonyCircuitBreaker[int])

	for i := 0; i < 3; i++ {
		_, err := cb.Execute(func() (int, error) {
			return 0, errors.New("fail")
		})
		assert.Error(t, err)
	}

	assert.Equal(t, Open, cb.CurrentState())
	wg.Wait()

	mu.Lock()
	assert.Contains(t, transitions, Open)
	mu.Unlock()
}

func TestSonyBreaker_OpenToHalfOpen(t *testing.T) {
	cb := NewSonyBreaker[int](&BreakerOptions{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		OpenTimeout:      50 * time.Millisecond,
	}).(*SonyCircuitBreaker[int])

	_, err := cb.Execute(func() (int, error) {
		return 0, errors.New("fail")
	})
	assert.Error(t, err)
	assert.Equal(t, Open, cb.CurrentState())

	time.Sleep(60 * time.Millisecond)
	assert.Equal(t, HalfOpen, cb.CurrentState())
}

func TestSonyBreaker_HalfOpenToClosed(t *testing.T) {
	cb := NewSonyBreaker[int](&BreakerOptions{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		OpenTimeout:      50 * time.Millisecond,
	}).(*SonyCircuitBreaker[int])

	_, _ = cb.Execute(func() (int, error) {
		return 0, errors.New("fail")
	})

	time.Sleep(60 * time.Millisecond)

	for i := 0; i < 2; i++ {
		val, err := cb.Execute(func() (int, error) {
			return 42, nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	}

	assert.Equal(t, Closed, cb.CurrentState())
}

func TestSonyBreaker_OpenRejectsCalls(t *testing.T) {
	cb := NewSonyBreaker[int](&BreakerOptions{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		OpenTimeout:      5 * time.Second,
	}).(*SonyCircuitBreaker[int])

	_, _ = cb.Execute(func() (int, error) {
		return 0, errors.New("fail")
	})

	val, err := cb.Execute(func() (int, error) {
		return 123, nil
	})
	assert.EqualError(t, err, "circuit breaker is open")
	assert.Equal(t, 0, val)
}

func TestSonyBreaker_OnStateChange(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	var mu sync.Mutex
	var calls []string

	cb := NewSonyBreaker[int](&BreakerOptions{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		OpenTimeout:      50 * time.Millisecond,
		OnStateChange: func(old, new State) {
			mu.Lock()
			calls = append(calls, old.String()+"-"+new.String())
			mu.Unlock()
			wg.Done()
		},
	}).(*SonyCircuitBreaker[int])

	_, _ = cb.Execute(func() (int, error) {
		return 0, errors.New("fail")
	})

	time.Sleep(60 * time.Millisecond)
	_, _ = cb.Execute(func() (int, error) {
		return 42, nil
	})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for OnStateChange, got: %v", calls)
	}

	mu.Lock()
	defer mu.Unlock()

	assert.Contains(t, calls, "closed-open")
	assert.Contains(t, calls, "open-half-open")
	assert.Contains(t, calls, "half-open-closed")
}
