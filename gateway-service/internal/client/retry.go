package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryConfig returns the default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// RetryableError determines if a gRPC error is retryable
func RetryableError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		return true // Non-gRPC errors are generally retryable
	}
	switch st.Code() {
	case codes.DeadlineExceeded,
		codes.Unavailable,
		codes.ResourceExhausted,
		codes.Aborted:
		return true
	default:
		return false
	}
}

// ExecuteWithRetry executes a gRPC call with exponential backoff retry
func ExecuteWithRetry(ctx context.Context, config RetryConfig, operation func(context.Context) error) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Check if context is cancelled before retrying
			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			default:
			}

			// Wait with exponential backoff
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return fmt.Errorf("context cancelled during retry: %w", ctx.Err())
			case <-timer.C:
			}

			// Increase delay for next attempt (exponential backoff)
			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}

		err := operation(ctx)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !RetryableError(err) {
			return err // Non-retryable error, return immediately
		}
	}

	return fmt.Errorf("max retries (%d) exceeded: %w", config.MaxRetries, lastErr)
}

// CallOption is a functional option for gRPC calls
type CallOption func(*callConfig)

type callConfig struct {
	retryConfig RetryConfig
	useRetry    bool
}

// WithRetry enables retry with the default configuration
func WithRetry() CallOption {
	return func(c *callConfig) {
		c.useRetry = true
		c.retryConfig = DefaultRetryConfig()
	}
}

// WithRetryConfig enables retry with custom configuration
func WithRetryConfig(config RetryConfig) CallOption {
	return func(c *callConfig) {
		c.useRetry = true
		c.retryConfig = config
	}
}

// GRPCInterceptor creates a client interceptor that adds retry logic
func GRPCInterceptor(config RetryConfig) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return ExecuteWithRetry(ctx, config, func(ctx context.Context) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		})
	}
}
