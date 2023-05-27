package grpc

import (
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"
	"user_service/grpc/service"
	"user_service/pkg/logger"
	"user_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvc client.ServiceManagerI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	user_service.RegisterUserServiceServer(grpcServer, service.NewUserService(cfg, log, strg, srvc))

	reflection.Register(grpcServer)
	return
}
