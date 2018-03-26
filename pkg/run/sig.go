package run

import (
	"context"
	"os"
	"os/signal"
)

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

		select {
		case <-sig:
			mainCanceler()
		case <-dismissc:
		case <-ctx.Done():
		}
	}()

	return ctx, dismisser
}
