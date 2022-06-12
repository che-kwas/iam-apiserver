package apiserver

import (
	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/server"
	"google.golang.org/grpc/reflection"

	pb "iam-apiserver/api/apiserver/proto/v1"
	"iam-apiserver/internal/apiserver/controller/cache"
	"iam-apiserver/internal/apiserver/store"
	"iam-apiserver/internal/apiserver/store/mysql"
)

type apiServer struct {
	*server.Server
	name string
	log  *logger.Logger

	err error
}

// NewServer builds a new apiServer.
func NewServer(name string) *apiServer {
	s := &apiServer{
		name: name,
		log:  logger.L(),
	}

	return s.initStore().newServer().setupHTTP().setupGRPC()
}

// Run runs the apiServer.
func (s *apiServer) Run() {
	defer s.log.Sync()
	if cli := store.Client(); cli != nil {
		defer cli.Close()
	}

	if s.err != nil {
		s.log.Fatal("failed to build the server: ", s.err)
	}

	if err := s.Server.Run(); err != nil {
		s.log.Fatal("server stopped unexpectedly: ", err)
	}
}

func (s *apiServer) initStore() *apiServer {
	var storeIns store.Store
	if storeIns, s.err = mysql.MySQLStore(); s.err != nil {
		return s
	}
	store.SetClient(storeIns)

	return s
}

func (s *apiServer) newServer() *apiServer {
	if s.err != nil {
		return s
	}

	s.Server, s.err = server.NewServer(s.name, server.WithGRPC())
	return s
}

func (s *apiServer) setupHTTP() *apiServer {
	if s.err != nil {
		return s
	}

	initRouter(s.Server.HTTPServer.Engine)
	return s
}

func (s *apiServer) setupGRPC() *apiServer {
	if s.err != nil {
		return s
	}

	ctrl := cache.NewCacheController()
	pb.RegisterCacheServer(s.GRPCServer, ctrl)
	reflection.Register(s.GRPCServer)

	return s
}
