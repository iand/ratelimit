// ratelimit implements rate limiting of functions using a leacky bucket algorithm
package ratelimit

import (
	"time"
)

type Func func()

type RateLimiter struct {
	ticker   *time.Ticker
	work     chan Func
	quit     chan struct{}
	draining bool
}

// Stop turns off the rate limiter immediately, losing any queued work
func (rl *RateLimiter) Stop() {
	close(rl.quit)
}

// Drain runs remaining work to completion and then stops the rate limiter
func (rl *RateLimiter) Drain() {
	rl.draining = true

	// block until we have drained
	<-rl.quit
}

// Do attempts to queue work for the rate limiter, returns false if it could not be queued
func (rl *RateLimiter) Do(fn Func) bool {
	if rl.draining {
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
				fn()
			default:
				// No work to do
				if rl.draining {
					close(rl.quit)
				}
			}
		}

	}
}

// PerSecond creates a ratelimiter that executes a maximum number of operations per second
func PerSecond(rate float64, capacity int) *RateLimiter {
	interval := time.Duration(float64(time.Second) * 1.0 / rate)
	rl := &RateLimiter{
		ticker: time.NewTicker(interval),
		quit:   make(chan struct{}),
		work:   make(chan Func, capacity),
	}

	go rl.run()
	return rl
}
