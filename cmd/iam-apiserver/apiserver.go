// iam-apiserver manages the api objects including users, policies, secrets.
package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/che-kwas/iam-kit/config"
	"github.com/spf13/pflag"

	"iam-apiserver/internal/apiserver"
)

var (
	name = "iam-apiserver"
	cfg  = pflag.StringP("config", "c", "", "config file")
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

	if err := config.LoadConfig(*cfg, name); err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	apiserver.NewServer(name).Run()
}
