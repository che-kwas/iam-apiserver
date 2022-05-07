package apiserver

import "github.com/che-kwas/iam-kit/server"

func NewServer(name string, cfgFile string) (*server.Server, error) {
	s, err := server.NewServer(name, cfgFile)
	if err != nil {
		return nil, err
	}

	s.InitRouter(initRouter)
	return s, nil
}
