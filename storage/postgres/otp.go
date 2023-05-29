package postgres

import (
	"database/sql"
	"fmt"
	"Projects/store/user_service/genproto/user_service"
	"Projects/store/user_service/pkg/helper"

	"context"

	"github.com/google/uuid"
)

func (u *userRepo) GetByPhoneNumber(ctx context.Context, req *user_service.UserPhoneNumberReq) (resp *user_service.User, err error) {

	query := `
		SELECT 
			id
		FROM "user"
		WHERE phone_number = $1
	`

	var (
		id sql.NullString
	)

	err = u.db.QueryRow(ctx, query, req.PhoneNumber).Scan(
		&id,
	)
	if err != nil {
		return resp, err
	}

	resp = &user_service.User{
		Id: id.String,
	}

	return
}

func (u *userRepo) CreateOTP(ctx context.Context, req *user_service.CreateOTP) error {

	id := uuid.New().String()

	query := `
		INSERT INTO "otp" (
			id,
			phone_number,
			otp
		) VALUES ($1, $2, $3)
	`

	otp, err := helper.GenerateOTP(6)
	if err != nil {
		return err
	}

	_, err = u.db.Exec(
		ctx,
		query,
		id,
		req.PhoneNumber,
		otp,
	)
	if err != nil {
		return err
	}

	fmt.Println(otp)

	return nil
}

func (u *userRepo) VerifyOTP(ctx context.Context, req *user_service.VerifyOTP) error {

	query := `
		SELECT 
			id,
			phone_number,
			otp,
			is_verify,
			created_at
		FROM "otp"
		WHERE phone_number = $1 AND otp = $2 AND is_verify = FALSE
		ORDER BY created_at DESC
	`

	var (
		id           sql.NullString
		phone_number sql.NullString
		otp          sql.NullString
		is_verify    sql.NullBool
		created_at   sql.NullString
	)

	err := u.db.QueryRow(
		ctx,
		query,
		req.PhoneNumber,
		req.Code,
	).Scan(&id, &phone_number, &otp, &is_verify, &created_at)
	if err != nil {
		return err
	}

	query = `
		UPDATE
			"otp"
		SET
			is_verify = TRUE
		WHERE id = $1
	`

	_, err = u.db.Exec(ctx, query, id.String)
	if err != nil {
		return err
	}

	return nil
}
