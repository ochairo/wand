// Package resilience provides resilience patterns like retry and circuit breaker.
package resilience

import (
	"fmt"
	"math"
	"time"
)

// RetryConfig configures the retry behavior with exponential backoff.
type RetryConfig struct {
	MaxAttempts       int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
}

// DefaultRetryConfig returns a retry configuration with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:       3,
		InitialDelay:      100 * time.Millisecond,
		MaxDelay:          10 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// Retrier retries operations with exponential backoff.
type Retrier struct {
	config RetryConfig
}

// NewRetrier creates a new Retrier with the specified configuration.
func NewRetrier(config RetryConfig) *Retrier {
	return &Retrier{config: config}
}

// Do executes the function with retry logic, returning the first successful result or the last error.
func (r *Retrier) Do(fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < r.config.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err
		if attempt < r.config.MaxAttempts-1 {
			delay := r.calculateBackoffDelay(attempt)
			time.Sleep(delay)
		}
	}
	return fmt.Errorf("all %d attempts failed: %w", r.config.MaxAttempts, lastErr)
}

func (r *Retrier) calculateBackoffDelay(attempt int) time.Duration {
	exponent := math.Min(float64(attempt), 10)
	delay := time.Duration(float64(r.config.InitialDelay) * math.Pow(r.config.BackoffMultiplier, exponent))
	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}
	return delay
}

// CircuitBreaker implements the circuit breaker pattern to prevent cascading failures.
type CircuitBreaker struct {
	state        string
	failureCount int
	threshold    int
	timeout      time.Duration
	lastFailTime time.Time
}

// NewCircuitBreaker creates a new CircuitBreaker with the specified threshold and timeout.
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{state: "closed", threshold: threshold, timeout: timeout}
}

// IsOpen returns true if the circuit is currently open (rejecting requests).
func (cb *CircuitBreaker) IsOpen() bool {
	if cb.state == "open" && time.Since(cb.lastFailTime) > cb.timeout {
		cb.state = "half-open"
		cb.failureCount = 0
		return false
	}
	return cb.state == "open"
}

// RecordSuccess resets the circuit to closed state after a successful operation.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.failureCount = 0
	cb.state = "closed"
}

// RecordFailure increments the failure count and opens the circuit if threshold is reached.
func (cb *CircuitBreaker) RecordFailure() {
	cb.failureCount++
	cb.lastFailTime = time.Now()
	if cb.failureCount >= cb.threshold {
		cb.state = "open"
	}
}

// State returns the current state of the circuit breaker ("open", "closed", or "half-open").
func (cb *CircuitBreaker) State() string {
	return cb.state
}
