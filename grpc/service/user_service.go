package service

import (
	"Projects/store/user_service/config"
	"Projects/store/user_service/genproto/user_service"
	"Projects/store/user_service/grpc/client"
	"Projects/store/user_service/models"
	"Projects/store/user_service/pkg/logger"
	"Projects/store/user_service/storage"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (u *UserService) Create(ctx context.Context, req *user_service.CreateUser) (resp *user_service.User, err error) {
	u.log.Info("---create user---", logger.Any("req", req))

	pkey, err := u.strg.User().Create(ctx, req)
	if err != nil {
		return nil, err
	}

	resp, err = u.strg.User().GetByPK(ctx, pkey)
	if err != nil {
		u.log.Error("!!!GetByPKeyUser->User->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, err
}

func (u *UserService) GetByID(ctx context.Context, req *user_service.UserPrimaryKey) (resp *user_service.User, err error) {

	u.log.Info("---get userbyid---", logger.Any("req", req))

	user, err := u.strg.User().GetByPK(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *UserService) GetList(ctx context.Context, req *user_service.GetListUserRequest) (resp *user_service.GetListUserResponse, err error) {

	u.log.Info("---GetUsers------>", logger.Any("req", req))

	resp, err = u.strg.User().GetList(ctx, req)
	if err != nil {
		u.log.Error("!!!GetUsers->User->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, err
}

func (i *UserService) Update(ctx context.Context, req *user_service.UpdateUser) (resp *user_service.User, err error) {

	i.log.Info("---UpdateUser------>", logger.Any("req", req))

	rowsAffected, err := i.strg.User().Update(ctx, req)

	if err != nil {
		i.log.Error("!!!UpdateUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.User().GetByPK(ctx, &user_service.UserPrimaryKey{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetUser->User->Get--->", logger.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *UserService) UpdatePatch(ctx context.Context, req *user_service.UpdatePatchUser) (resp *user_service.User, err error) {

	i.log.Info("---UpdatePatchUser------>", logger.Any("req", req))

	updatePatchModel := models.UpdatePatchRequest{
		Id:     req.GetId(),
		Fields: req.GetFields().AsMap(),
	}

	rowsAffected, err := i.strg.User().UpdatePatch(ctx, &updatePatchModel)

	if err != nil {
		i.log.Error("!!!UpdatePatchUser--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if rowsAffected <= 0 {
		return nil, status.Error(codes.InvalidArgument, "no rows were affected")
	}

	resp, err = i.strg.User().GetByPK(ctx, &user_service.UserPrimaryKey{Id: req.Id})
	if err != nil {
		i.log.Error("!!!GetUser->User->Get--->", logger.Error(err))

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return resp, err
}

func (i *UserService) Delete(ctx context.Context, req *user_service.UserPrimaryKey) (resp *empty.Empty, err error) {

	i.log.Info("---DeleteUser------>", logger.Any("req", req))

	err = i.strg.User().Delete(ctx, req)
	if err != nil {
		i.log.Error("!!!DeleteUser->User->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}

// otp
func (i *UserService) CreateUserOTP(ctx context.Context, req *user_service.CreateOTP) (resp *empty.Empty, err error) {

	i.log.Info("---CreateUserOTP------->", logger.Any("req", req))

	err = i.strg.User().CreateOTP(ctx, req)
	if err != nil {
		i.log.Error("!!!CreateUserOTP->OTP->Create--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}

func (i *UserService) VerifyUserOTP(ctx context.Context, req *user_service.VerifyOTP) (resp *empty.Empty, err error) {

	i.log.Info("---VerifyUserOTP------->", logger.Any("req", req))

	err = i.strg.User().VerifyOTP(ctx, req)
	if err != nil {
		i.log.Error("!!!VerifyUserOTP->OTP->Verify--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &empty.Empty{}, nil
}

func (i *UserService) Check(ctx context.Context, req *user_service.UserPhoneNumberReq) (resp *user_service.User, err error) {

	i.log.Info("---CheckUser------>", logger.Any("req", req))

	resp, err = i.strg.User().GetByPhoneNumber(ctx, req)
	if err != nil {
		i.log.Error("!!!GetUserByPhoneNumber->User->Get--->", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return resp, err
}
