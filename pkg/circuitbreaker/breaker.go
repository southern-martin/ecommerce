package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

// Config holds the configuration for a circuit breaker instance.
type Config struct {
	Name             string        // Name identifying this circuit breaker.
	MaxRequests      uint32        // Max requests allowed in half-open state (default 3).
	Interval         time.Duration // Closed state counters reset interval (default 60s).
	Timeout          time.Duration // Duration the breaker stays open before transitioning to half-open (default 30s).
	FailureThreshold int           // Consecutive failures required to trip the breaker (default 5).
}

// defaults fills in zero-valued fields with sensible defaults.
func (cfg *Config) defaults() {
	if cfg.MaxRequests == 0 {
		cfg.MaxRequests = 3
	}
	if cfg.Interval == 0 {
		cfg.Interval = 60 * time.Second
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
}

// Breaker wraps sony/gobreaker.CircuitBreaker providing a simplified API.
// The breaker transitions through three states:
//   - Closed: normal operation; requests pass through.
//   - Open: requests are immediately rejected; the breaker waits for Timeout.
//   - HalfOpen: a limited number of requests (MaxRequests) are allowed through
//     to test if the downstream service has recovered.
type Breaker struct {
	cb *gobreaker.CircuitBreaker
}

// New creates a new Breaker with the given configuration. Zero-valued fields
// in cfg are replaced with defaults.
func New(cfg Config) *Breaker {
	cfg.defaults()

	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip when consecutive failures exceed the threshold.
			if int(counts.ConsecutiveFailures) > cfg.FailureThreshold {
				return true
			}
			// Also trip when the failure ratio exceeds 60% with a minimum
			// number of requests to avoid tripping on initial sparse traffic.
			if counts.Requests >= 10 {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return failureRatio > 0.6
			}
			return false
		},
	}

	return &Breaker{
		cb: gobreaker.NewCircuitBreaker(settings),
	}
}

// Execute runs the given function through the circuit breaker. If the breaker
// is open, Execute returns immediately with gobreaker.ErrOpenState. If the
// function returns an error, it is counted as a failure.
func (b *Breaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return b.cb.Execute(fn)
}

// State returns the current state of the circuit breaker as a human-readable
// string: "closed", "half-open", or "open".
func (b *Breaker) State() string {
	switch b.cb.State() {
	case gobreaker.StateClosed:
		return "closed"
	case gobreaker.StateHalfOpen:
		return "half-open"
	case gobreaker.StateOpen:
		return "open"
	default:
		return "unknown"
	}
}
