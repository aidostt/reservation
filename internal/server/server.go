package server

import (
	"dip/internal/grpcauth"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	GrpcServer *grpc.Server
}

func NewServer() *Server {
	return &Server{
		GrpcServer: grpc.NewServer(grpc.UnaryInterceptor(grpcauth.UnaryInterceptor)),
	}

}

// Run runs gRPC server.
func (s *Server) Run(l net.Listener) error {
	return s.GrpcServer.Serve(l)

}

// Stop stops gRPC server.
func (s *Server) Stop() {
	s.GrpcServer.GracefulStop()
}
