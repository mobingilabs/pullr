package run

import (
	"context"
	"time"
)

// Retry runs the given job til it succeed in given duration. It also waits for
// a small amount of time defined by the delay parameter between retries. It
// returns the last error if timeout duration passed otherwise returns nil.
func Retry(timeout, delay time.Duration, job func() error) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return RetryWithContext(ctx, delay, job)
}

// RetryWithContext runs the given job til it succeed in given duration. If the
// given context done somehow it will return context's error. It also waits for
// a small amount of time defined by the delay parameter between retries. It
// returns the last error if timeout duration passed otherwise returns nil.
func RetryWithContext(ctx context.Context, delay time.Duration, job func() error) (err error) {
	donec := make(chan struct{})
	ticker := time.NewTicker(delay)
	go func() {
		defer func() { recover() }()
		// We want the first tick without waiting delay
		for ok := true; ok; _, ok = <-ticker.C {
			err = job()
			if err == nil {
				donec <- struct{}{}
				break
			}
		}
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-donec:
		err = nil
	}

	ticker.Stop()
	close(donec)
	return err
}
