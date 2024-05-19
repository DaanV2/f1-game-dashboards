package api

import (
	"errors"
	"fmt"
	"net"

	grpc_gen "github.com/DaanV2/f1-game-dashboards/server/api/grpc"
	"github.com/DaanV2/f1-game-dashboards/server/sessions"
	"github.com/charmbracelet/log"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServerOptions struct {
	port string
	host string
}

type grpcServer struct {
	grpc_gen.UnimplementedChairServiceServer

	chairs *sessions.ChairManager

	grpc *grpc.Server

	options grpcServerOptions
}

func newGrpcServer(chairs *sessions.ChairManager, options grpcServerOptions) *grpcServer {
	return &grpcServer{
		UnimplementedChairServiceServer: grpc_gen.UnimplementedChairServiceServer{},
		chairs:                          chairs,
		options:                         options,
		grpc:                            nil,
	}
}

func (s *grpcServer) Start() error {
	address := fmt.Sprintf("%s:%s", s.options.host, s.options.port)
	log.Info("starting grpc server...", "address", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption

	// TODO add health checking

	s.grpc = grpc.NewServer(opts...)
	grpc_gen.RegisterChairServiceServer(s.grpc, s)
	reflection.Register(s.grpc)

	go func() {
		
		err := s.grpc.Serve(lis)
		if errors.Is(err, grpc.ErrServerStopped) {
			log.Info("grpc server stopped")
		} else {
			log.Error("grpc server stopped with error", "error", err)
		}
	}()

	return nil
}

func (s *grpcServer) Stop() error {
	log.Info("stopping grpc server...")
	defer log.Info("grpc server stopped")
	s.grpc.GracefulStop()

	return nil
}