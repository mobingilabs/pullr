package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mobingilabs/pullr/cmd/apisrv/app"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/sirupsen/logrus"
)

func main() {
	pemCerts, err := ioutil.ReadFile("./ca-certificates.crt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read ca-certificates: %v", err)
		os.Exit(1)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)

	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}

	errs.SetLogger(logrus.StandardLogger())
	errs.Fatal(app.RootCmd.Execute())
}
