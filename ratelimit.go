// Package ratelimit implements rate limiting of functions using a leaky bucket algorithm
package ratelimit

import (
	"context"
	"sync/atomic"
	"time"
)

type Func func()

type RateLimiter struct {
	ticker   *time.Ticker
	work     chan Func
	quit     chan struct{}
	draining uint64
}

// Stop turns off the rate limiter immediately, losing any queued work
func (rl *RateLimiter) Stop() {
	close(rl.quit)
}

// Drain runs remaining work to completion and then stops the rate limiter
func (rl *RateLimiter) Drain() {
	atomic.StoreUint64(&rl.draining, 1)

	// block until we have drained
	<-rl.quit
}

// Do attempts to queue work for the rate limiter, returns false if it could not be queued.
// Each function queued will be executed in a separate goroutine so if the functions are long running, this
// could result in a large number of active goroutines, depending on the rate of the limiter
func (rl *RateLimiter) Do(ctx context.Context, fn Func) bool {
	if ctx != nil {
		// Check whether context has been cancelled
		select {
		case <-ctx.Done():
			return false
		default:
		}
	}

	if draining := atomic.LoadUint64(&rl.draining); draining == 1 {
		return false
	}
	select {
	case rl.work <- fn:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) run() {
	for {
		select {
		case <-rl.quit:
			rl.ticker.Stop()
			return
		case <-rl.ticker.C:
			select {
			case fn := <-rl.work:
				go fn()
			default:
				// No work to do
				if draining := atomic.LoadUint64(&rl.draining); draining == 1 {
					close(rl.quit)
				}
			}
		}

	}
}

// PerSecond creates a ratelimiter that executes a maximum number of operations per second
func PerSecond(rate float64, capacity int) *RateLimiter {
	interval := time.Duration(float64(time.Second) / rate)
	rl := &RateLimiter{
		ticker: time.NewTicker(interval),
		quit:   make(chan struct{}),
		work:   make(chan Func, capacity),
	}

	go rl.run()
	return rl
}
