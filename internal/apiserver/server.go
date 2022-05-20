package apiserver

import (
	"log"

	"github.com/che-kwas/iam-kit/config"
	"github.com/che-kwas/iam-kit/server"

	"iam-apiserver/internal/apiserver/store"
	"iam-apiserver/internal/apiserver/store/mysql"
)

type apiServer struct {
	*server.Server
	name string
	err  error
}

// NewServer builds a new apiServer.
func NewServer(name, cfgFile string) *apiServer {
	s := &apiServer{name: name}
	return s.loadConfig(cfgFile).initStore().initCache().build()
}

// Run runs the apiServer.
func (s *apiServer) Run() {
	if s.err != nil {
		log.Fatal("Build server error: ", s.err)
	}

	if err := s.Server.Run(); err != nil {
		log.Fatal("Server stopped unexpectedly: ", err)
	}
}

func (s *apiServer) loadConfig(cfgFile string) *apiServer {
	s.err = config.LoadConfig(cfgFile, s.name)
	return s
}

func (s *apiServer) initStore() *apiServer {
	if s.err != nil {
		return s
	}

	var storeIns store.Store
	storeIns, s.err = mysql.GetMySQLStore()
	if s.err != nil {
		return s
	}

	store.SetClient(storeIns)
	return s
}

func (s *apiServer) initCache() *apiServer {
	if s.err != nil {
		return s
	}

	return s
}

func (s *apiServer) build() *apiServer {
	if s.err != nil {
		return s
	}

	s.Server, s.err = server.NewServer(s.name)
	if s.err != nil {
		return s
	}

	s.InitRouter(initRouter)
	return s
}
