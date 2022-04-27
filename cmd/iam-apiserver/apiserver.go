// The IAM API server manages the api objects including users, policies, secrets and more.

package main

import (
	log "github.com/sirupsen/logrus"

	"iam-apiserver/internal/apiserver"
)

func main() {
	// TODO get cfgFile from flags
	cfgFile := ""

	server, err := apiserver.NewServer("iam-apiserver", cfgFile)
	if err != nil {
		log.Fatal("Build server error: ", err)
	}

	if err := server.Run(); err != nil {
		log.Fatal("Server stopped unexpectedly: ", err)
	}
}
