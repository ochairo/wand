package resilience

import (
	"fmt"
	"testing"
	"time"
)

func TestRetrySuccess(t *testing.T) {
	retrier := NewRetrier(DefaultRetryConfig())
	attempts := 0
	err := retrier.Do(func() error {
		attempts++
		if attempts < 2 {
			return fmt.Errorf("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected success, got: %v", err)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempts)
	}
}

func TestRetryMaxAttempts(t *testing.T) {
	config := DefaultRetryConfig()
	config.MaxAttempts = 3
	retrier := NewRetrier(config)

	attempts := 0
	err := retrier.Do(func() error {
		attempts++
		return fmt.Errorf("persistent error")
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestCircuitBreakerOpen(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	if cb.IsOpen() {
		t.Error("Circuit should be closed initially")
	}

	cb.RecordFailure()
	if cb.IsOpen() {
		t.Error("Circuit should not be open after 1 failure")
	}

	cb.RecordFailure()
	if !cb.IsOpen() {
		t.Error("Circuit should be open after 2 failures")
	}
}

func TestCircuitBreakerReset(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	cb.RecordFailure()
	cb.RecordFailure()
	if !cb.IsOpen() {
		t.Error("Circuit should be open")
	}

	cb.RecordSuccess()
	if cb.IsOpen() {
		t.Error("Circuit should be closed after success")
	}
}
