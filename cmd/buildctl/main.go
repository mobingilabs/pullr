package main

import (
	"github.com/mobingilabs/pullr/cmd/buildctl/app"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/sirupsen/logrus"
)

func main() {
	errs.SetLogger(logrus.StandardLogger())
	errs.Fatal(app.RootCmd.Execute())
}
