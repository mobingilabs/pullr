package main

import (
	"github.com/mobingilabs/pullr/cmd/apisrv/app"
	"github.com/mobingilabs/pullr/pkg/errs"
)

func main() {
	errs.Fatal(app.Execute())
}
