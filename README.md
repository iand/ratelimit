ratelimit
=========

Go package for rate limiting of functions using a leaky bucket algorithm

[![Build Status](https://travis-ci.org/iand/pgen.svg?branch=master)](https://travis-ci.org/iand/pgen)

## Installation

Simply run

    go get -u github.com/iand/ratelimit

Documentation is at [http://godoc.org/github.com/iand/ratelimit](http://godoc.org/github.com/iand/ratelimit)

## Usage

Example use:

```go
import "github.com/iand/ratelimit"

func main() {
	limiter := ratelimit.PerSecond(5, 100)

	// Following should be executed at a rate of 5 per second
	for i := 0; i < 100; i++ {
		limiter.Do(func() { println("hello") })
	}

	// Blocks until the rate limiter has finished
	limiter.Drain()
}
```

LICENSE
=======
This is free and unencumbered software released into the public domain. For more
information, see <http://unlicense.org/> or the accompanying [`UNLICENSE`](UNLICENSE) file.
