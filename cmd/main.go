package main

import (
	"context"
	"net"
	"user_service/grpc"
	"user_service/grpc/client"
	"user_service/pkg/logger"
	"user_service/storage/postgres"
	"user_service/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	loggerLevel := logger.LevelDebug

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	pgStore, err := postgres.NewPostgres(context.Background(), cfg)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}
	defer pgStore.CloseDB()

	svcs, err := client.NewGrpcClients(cfg)
	if err != nil {
		log.Panic("client.NewGrpcClients", logger.Error(err))
	}

	grpcServer := grpc.SetUpServer(cfg, log, pgStore, svcs)

	lis, err := net.Listen("tcp", cfg.ServicePort)
	if err != nil {
		log.Panic("net.Listen", logger.Error(err))
	}

	log.Info("GRPC: Server being started...", logger.String("port", cfg.ServicePort))

	if err := grpcServer.Serve(lis); err != nil {
		log.Panic("grpcServer.Serve", logger.Error(err))
	}
}
