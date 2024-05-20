package api

import (
	"errors"

	"github.com/DaanV2/f1-game-dashboards/server/authenication"
	"github.com/DaanV2/f1-game-dashboards/server/sessions"
)

type apiServerOptions struct {
	grpc grpcServerOptions
}

type ApiServer struct {
	options apiServerOptions

	grpcServer *grpcServer
}

func NewApiServer(chairs *sessions.ChairManager, authenicator *authenication.Authenticator) *ApiServer {
	options := apiServerOptions{
		grpc: grpcServerOptions{
			port: "50051",
			host: "0.0.0.0",
		},
	}

	return &ApiServer{
		options: options,

		grpcServer: newGrpcServer(chairs, authenicator, options.grpc),
	}
}

func (server *ApiServer) Start() error {
	return errors.Join(
		server.grpcServer.Start(),
	)
}

func (server *ApiServer) Stop() error {
	return errors.Join(
		server.grpcServer.Stop(),
	)
}
