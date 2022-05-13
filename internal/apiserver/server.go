package apiserver

import (
	"iam-apiserver/internal/apiserver/store"
	"iam-apiserver/internal/apiserver/store/mysql"

	"github.com/che-kwas/iam-kit/config"
	"github.com/che-kwas/iam-kit/server"
)

func NewServer(name string, cfgFile string) (*server.Server, error) {
	if err := config.LoadConfig(cfgFile, name); err != nil {
		return nil, err
	}

	storeIns, err := mysql.GetMySQLStore()
	if err != nil {
		return nil, err
	}
	store.SetClient(storeIns)

	s, err := server.NewServer(name)
	if err != nil {
		return nil, err
	}

	s.InitRouter(initRouter)
	return s, nil
}
