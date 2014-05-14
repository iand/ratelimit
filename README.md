ratelimit
=========

Go package for rate limiting of functions using a leaky bucket algorithm

Docs: http://godoc.org/github.com/iand/ratelimit

Example use:

```go
import "github.com/iand/ratelimit"

limiter := ratelimit.PerSecond(5, 100)

// Following should be executed at a rate of 5 per seconf
for i := 0; i < 100; i++ {
	limiter.Do(func() { println("hello") })
}

// Blocks until the rate limiter has finished
limiter.Drain()
```
