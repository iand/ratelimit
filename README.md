ratelimit
=========

Go package for rate limiting of functions using a leacky bucket algorithm

Docs: http://godoc.org/github.com/iand/ratelimit

Example use:

```
import "github.com/iand/ratelimit"

limiter := ratelimit.PerSecond(5, 100)

// Following should be executed at a rate of 5 per seconf
for i := 0; i < 100; i++ {
	limiter.Do(func() { println("hello") })
}
```
