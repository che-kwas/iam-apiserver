// The IAM API server manages the api objects including users, policies, secrets and more.

package main

import (
	"log"

	"github.com/spf13/pflag"

	"iam-apiserver/internal/apiserver"
)

var (
	name = "iam-apiserver"
	cfg  = pflag.StringP("config", "c", "./iam-apiserver.yaml", "config file")
	help = pflag.BoolP("help", "h", false, "show help message")
)

func main() {
	pflag.Parse()
	if *help {
		pflag.Usage()
		return
	}

	server, err := apiserver.NewServer(name, *cfg)
	if err != nil {
		log.Fatal("Build server error: ", err)
	}

	if err := server.Run(); err != nil {
		log.Fatal("Server stopped unexpectedly: ", err)
	}
}
