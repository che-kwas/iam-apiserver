package apiserver

import (
	"log"

	"github.com/che-kwas/iam-kit/server"

	"iam-apiserver/internal/apiserver/store"
	"iam-apiserver/internal/apiserver/store/mysql"
)

type apiServer struct {
	*server.Server
	name string

	err error
}

// NewServer builds a new apiServer.
func NewServer(name string) *apiServer {
	s := &apiServer{name: name}

	return s.initStore().newServer().registerRouter()
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

func (s *apiServer) initStore() *apiServer {
	var storeIns store.Store
	storeIns, s.err = mysql.MySQLStore()
	if s.err != nil {
		return s
	}

	store.SetClient(storeIns)
	return s
}

func (s *apiServer) newServer() *apiServer {
	if s.err != nil {
		return s
	}

	s.Server, s.err = server.NewServer(s.name)
	return s
}

func (s *apiServer) registerRouter() *apiServer {
	if s.err != nil {
		return s
	}

	s.InitRouter(initRouter)
	return s
}
