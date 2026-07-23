package middleware

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	mu               sync.Mutex
	state            State
	failureCount     int
	successCount     int
	maxFailures      int
	successThreshold int
	timeout          time.Duration
	lastFailureTime  time.Time
}

func NewCircuitBreaker(maxFailures, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		maxFailures:      maxFailures,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.failureCount = 0
			cb.successCount = 0
		} else {
			return errors.New("circuit breaker is open")
		}
	case StateHalfOpen:
	}

	err := fn()
	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		if cb.state == StateHalfOpen || cb.failureCount >= cb.maxFailures {
			cb.state = StateOpen
		}
		return err
	}

	if cb.state == StateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.state = StateClosed
			cb.failureCount = 0
			cb.successCount = 0
		}
	} else {
		cb.failureCount = 0
	}

	return nil
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
	cb.successCount = 0
}

