package ratelimit

import (
	"testing"
	"time"
)

func TestPerSecond(t *testing.T) {

	expectedRate := 500.0
	sampleSize := 1000

	done := make(chan struct{}, sampleSize)
	limiter := PerSecond(expectedRate, sampleSize)

	start := time.Now()
	for i := 0; i < sampleSize; i++ {
		limiter.Do(func() {
			done <- struct{}{}
		})
	}
	limiter.Drain()

	for i := 0; i < sampleSize; i++ {
		<-done
	}

	actualRate := float64(sampleSize) / time.Now().Sub(start).Seconds()

	if actualRate > expectedRate {
		t.Errorf("Rate of %0.4f exceeded expected rate of %0.4f", actualRate, expectedRate)
	}

}
