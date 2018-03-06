package app

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/mobingilabs/pullr/cmd/trackersvc/service"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ListenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Start listening for status updates",
		Long:  "Start listening for status updates",
		Run: func(cmd *cobra.Command, args []string) {
			rand.Seed(time.Now().UnixNano())

			mainCtx, sigCanceler := errs.ContextWithSig(context.Background(), os.Interrupt, os.Kill)
			defer sigCanceler()

			timeoutCtx, timeoutCanceler := context.WithTimeout(mainCtx, time.Minute*5)
			tracker, err := service.New(timeoutCtx, Logger, Config)
			timeoutCanceler()
			if err != nil {
				if errors.Cause(err) == context.Canceled {
					logrus.Info("Program interrupted! Terminated gracefully.")
					return
				}
				logrus.Fatalf("failed to init Tracker service: %v", err)
			}

			if err := tracker.Listen(mainCtx); err != nil {
				if errors.Cause(err) == context.Canceled {
					logrus.Info("Program interrupted! Terminated gracefully.")
					return
				}
				logrus.Fatalf("Tracker crashed with: %v", err)
			}
		},
	}
)
