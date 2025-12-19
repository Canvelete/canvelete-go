package canvelete

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxAttempts   int
	BackoffFactor float64
	InitialDelay  time.Duration
	MaxDelay      time.Duration
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   3,
		BackoffFactor: 2.0,
		InitialDelay:  time.Second,
		MaxDelay:      60 * time.Second,
	}
}

// RetryableFunc is a function that can be retried
type RetryableFunc func() error

// WithRetry executes a function with retry logic
func WithRetry(ctx context.Context, config *RetryConfig, fn RetryableFunc) error {
	if config == nil {
		config = DefaultRetryConfig()
	}
	
	var lastErr error
	delay := config.InitialDelay
	
	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		
		// Check if error is retryable
		if !isRetryableError(lastErr) {
			return lastErr
		}
		
		// Check for rate limit with retry-after
		if apiErr, ok := lastErr.(*APIError); ok && apiErr.StatusCode == http.StatusTooManyRequests {
			if retryAfter := getRetryAfter(apiErr); retryAfter > 0 {
				delay = retryAfter
			}
		}
		
		// Don't sleep on last attempt
		if attempt < config.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			
			// Exponential backoff
			delay = time.Duration(float64(delay) * config.BackoffFactor)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}
	}
	
	return lastErr
}

// isRetryableError checks if an error should be retried
func isRetryableError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		// Retry on rate limit and server errors
		return apiErr.StatusCode == http.StatusTooManyRequests ||
			apiErr.StatusCode >= http.StatusInternalServerError
	}
	return false
}

// getRetryAfter extracts retry-after duration from an API error
func getRetryAfter(err *APIError) time.Duration {
	// This would need to be enhanced to parse Retry-After header
	// For now, return 0 to use default backoff
	return 0
}

// RetryOnRateLimit is a convenience function for retrying on rate limits
func RetryOnRateLimit(ctx context.Context, fn RetryableFunc) error {
	config := &RetryConfig{
		MaxAttempts:   5,
		BackoffFactor: 2.0,
		InitialDelay:  time.Second,
		MaxDelay:      60 * time.Second,
	}
	return WithRetry(ctx, config, fn)
}

// parseRetryAfterHeader parses the Retry-After header value
func parseRetryAfterHeader(value string) time.Duration {
	// Try parsing as seconds
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	
	// Try parsing as HTTP date
	if t, err := http.ParseTime(value); err == nil {
		return time.Until(t)
	}
	
	return 0
}

// calculateBackoff calculates the backoff duration for a given attempt
func calculateBackoff(attempt int, config *RetryConfig) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt))
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}
	return time.Duration(delay)
}