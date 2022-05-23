// The IAM API server manages the api objects including users, policies, secrets and more.

package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/spf13/pflag"

	"iam-apiserver/internal/apiserver"
	"iam-apiserver/internal/pkg/config"
)

var (
	name = "iam-apiserver"
	cfg  = pflag.StringP("config", "c", "./iam-apiserver.yaml", "config file")
	help = pflag.BoolP("help", "h", false, "show help message")
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// parse flag
	pflag.Parse()
	if *help {
		pflag.Usage()
		return
	}

	if err := config.InitConfig(*cfg, name); err != nil {
		log.Fatal("Initializa config failed: ", err)
	}

	apiserver.NewServer(name).Run()
}
