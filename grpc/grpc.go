package grpc

import (
	"Projects/store/user_service/config"
	"Projects/store/user_service/genproto/user_service"
	"Projects/store/user_service/grpc/client"
	"Projects/store/user_service/grpc/service"
	"Projects/store/user_service/pkg/logger"
	"Projects/store/user_service/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvc client.ServiceManagerI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	user_service.RegisterUserServiceServer(grpcServer, service.NewUserService(cfg, log, strg, srvc))

	reflection.Register(grpcServer)
	return
}
