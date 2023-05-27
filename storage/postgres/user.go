package postgres

import (
	"context"
	"database/sql"
	"user_service/genproto/user_service"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, req *user_service.CreateUser) (resp *user_service.UserPrimaryKey, err error) {
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

	_, err = u.db.Exec(
		ctx,
		query,
		id,
		req.FirstName,
		req.LastName,
		req.PhoneNumber,
		req.DateOfBirth,
		req.Password,
	)
	if err != nil {
		return nil, err
	}

	return &user_service.UserPrimaryKey{Id: id}, nil
}

func (u *userRepo) GetByPKey(ctx context.Context, req *user_service.UserPrimaryKey) (user *user_service.User, err error) {

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			phone_number, 
			date_of_birth
			password,
			created_at,
			updated_at
		FROM "user"
		WHERE id = $1
	`

	var (
		id            sql.NullString
		first_name    sql.NullString
		last_name     sql.NullString
		phone_number  sql.NullString
		date_of_birth sql.NullString
		password      sql.NullString
		created_at    sql.NullString
		updated_at    sql.NullString
	)

	err = u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&first_name,
		&last_name,
		&phone_number,
		&date_of_birth,
		&password,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return user, err
	}

	user = &user_service.User{
		Id:          id.String,
		FirstName:   first_name.String,
		LastName:    last_name.String,
		PhoneNumber: phone_number.String,
		DateOfBirth: date_of_birth.String,
		Password: password.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
	}

	return user, nil
}
