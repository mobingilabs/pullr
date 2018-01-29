package main

import (
	"github.com/golang/glog"
	"github.com/mobingilabs/pullr/cmd/apisrv/app"
)

func main() {
	glog.CopyStandardLogTo("INFO")
	app.Execute()
}
