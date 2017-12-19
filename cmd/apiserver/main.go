package main

import (
	"github.com/golang/glog"
	"github.com/mobingilabs/pullr/cmd/apiserver/app"
)

func main() {
	glog.CopyStandardLogTo("INFO")
	app.Execute()
}
