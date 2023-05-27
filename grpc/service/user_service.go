package service

import (
	"context"
	"user_service/config"
	"user_service/genproto/user_service"
	"user_service/grpc/client"
	"user_service/pkg/logger"
	"user_service/storage"
)

type UserService struct {
	cfg      config.Config
	log      logger.LoggerI
	strg     storage.StorageI
	services client.ServiceManagerI
	*user_service.UnimplementedUserServiceServer
}

func NewUserService(cfg config.Config, log logger.LoggerI, strg storage.StorageI, srvs client.ServiceManagerI) *UserService {
	return &UserService{
		cfg:      cfg,
		log:      log,
		strg:     strg,
		services: srvs,
	}
}

func (u *UserService) Create(ctx context.Context, req *user_service.CreateUser) (*user_service.UserPrimaryKey, error) {
	u.log.Info("---create user---",logger.Any("req",req))

	pkey,err:=u.strg.User().Create(ctx,req)
	if err != nil {
		return nil,err
	}

	resp:=u.strg.User().Get
}
