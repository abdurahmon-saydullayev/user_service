package postgres

import (
	"context"
	"user_service/genproto/user_service"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
)

type usertRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *usertRepo {
	return &usertRepo{
		db: db,
	}
}

func (u *usertRepo) Create(ctx context.Context, req *user_service.CreateUser) (resp *user_service.UserPrimaryKey, err error) {
	id := uuid.New().String()

	query := `
	INSERT INTO user (
		id,
		first_name,
		last_name,
		phone_number,
		date_of_birth,
		password,
		created_at,
		updated_at,
	) VALUES($1,$2,$3,$4,$5,$6, NOW(),NOW())
	`

	_, err=u.db.Exec(
		ctx,
		query,
		id,
		req.FirstName,
		req.LastName,
		req.PhoneNumber,
		req.DateOfBirth,
		req.Password,
	)
	if err !=nil{
		return nil,err
	}

	return &user_service.UserPrimaryKey{Id: id}, nil
}
