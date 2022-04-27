package apiserver

import "github.com/che-kwas/iam-kit/server"

func NewServer(name string, cfgFile string) (*server.Server, error) {
	// TODO: add routers
	return server.NewServer(name, cfgFile)
}
