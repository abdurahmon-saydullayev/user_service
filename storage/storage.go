package storage

import (
	"context"
	"user_service/genproto/user_service"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
}

type UserRepoI interface {
	Create(context.Context, *user_service.CreateUser) (*user_service.UserPrimaryKey, error)
	// GetByPK(context.Context, *user_service.UserPrimaryKey) (*user_service.User, error)
	// GetList(context.Context, *user_service.GetListUserRequest) (*user_service.GetListUserResponse, error)
	// Update(context.Context, *user_service.UserPrimaryKey) (int64, error)
	// UpdatePatch(context.Context, *user_service.UserPrimaryKey) error

	// GetByPhoneNumber(context.Context, *user_service.UserPhoneNumberReq) (*user_service.User, error)

	// //otp
	// CreateOTP(context.Context, *user_service.CreateOTP) error
	// VerifyOTP(context.Context, *user_service.VerifyOTP) error
}
