package storage

import (
	"context"
	"Projects/store/user_service/genproto/user_service"
	"Projects/store/user_service/models"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
}

type UserRepoI interface {
	Create(context.Context, *user_service.CreateUser) (*user_service.UserPrimaryKey, error)
	GetByPK(context.Context, *user_service.UserPrimaryKey) (*user_service.User, error)
	GetList(context.Context, *user_service.GetListUserRequest) (*user_service.GetListUserResponse, error)
	Update(ctx context.Context, req *user_service.UpdateUser) (rowsAffected int64, err error)
	UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (rowsAffected int64, err error)
	Delete(ctx context.Context, req *user_service.UserPrimaryKey) error

	// GetByPhoneNumber(context.Context, *user_service.UserPhoneNumberReq) (*user_service.User, error)

	// //otp
	// CreateOTP(context.Context, *user_service.CreateOTP) error
	// VerifyOTP(context.Context, *user_service.VerifyOTP) error
}
