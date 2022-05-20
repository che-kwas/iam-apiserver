// The IAM API server manages the api objects including users, policies, secrets and more.

package main

import (
	"math/rand"
	"time"

	"github.com/spf13/pflag"

	"iam-apiserver/internal/apiserver"
)

var (
	name = "iam-apiserver"
	cfg  = pflag.StringP("config", "c", "./iam-apiserver.yaml", "config file")
	help = pflag.BoolP("help", "h", false, "show help message")
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	pflag.Parse()
	if *help {
		pflag.Usage()
		return
	}

	apiserver.NewServer(name, *cfg).Run()
}
