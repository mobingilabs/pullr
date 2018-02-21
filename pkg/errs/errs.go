package errs

import (
	"context"
	"os"
	"os/signal"
	"time"
)

// ErrLogger interface has error logging functions
type ErrLogger interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
}

var logger ErrLogger

// SetLogger sets errs package's global logger
func SetLogger(errLogger ErrLogger) {
	logger = errLogger
}

// Log checks if the given err is null and if not it logs it
func Log(err error) {
	if err != nil {
		logger.Error(err)
	}
}

// Fatal checks if the given err is null and if not it logs and exits the program
func Fatal(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

// Retry runs the given job til it succeed in given duration. It also waits for
// a small amount of time defined by the delay parameter between retries. It
// returns the last error if timeout duration passed otherwise returns nil.
func Retry(timeout, delay time.Duration, job func() error) (err error) {
	return RetryWithContext(context.Background(), timeout, delay, job)
}

// RetryWithContext runs the given job til it succeed in given duration. If the
// given context done somehow it will return context's error. It also waits for
// a small amount of time defined by the delay parameter between retries. It
// returns the last error if timeout duration passed otherwise returns nil.
func RetryWithContext(ctx context.Context, timeout, delay time.Duration, job func() error) (err error) {
	if timeout < delay {
		timeout = delay
	}

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
		ticker.Stop()
		close(donec)
		return ctx.Err()
	case <-time.After(timeout):
		ticker.Stop()
		close(donec)
		return err
	case <-donec:
		ticker.Stop()
		close(donec)
		return nil
	}
}

// ContextWithSig creates a new context with passed context as it's parent and a
// dismiss function. Dismiss function stops listening for os signals. Please
// make sure to call dismiss function.
func ContextWithSig(ctx context.Context, sigs ...os.Signal) (context.Context, func()) {
	ctx, mainCanceler := context.WithCancel(ctx)
	dismissc := make(chan struct{}, 1)
	dismisser := func() { dismissc <- struct{}{}; close(dismissc) }

	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, sigs...)

		defer signal.Stop(sig)
		defer close(sig)

		select {
		case <-sig:
			mainCanceler()
		case <-dismissc:
		case <-ctx.Done():
		}
	}()

	return ctx, dismisser
}
