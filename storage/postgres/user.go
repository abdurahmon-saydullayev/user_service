package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"Projects/store/user_service/genproto/user_service"
	"Projects/store/user_service/models"
	"Projects/store/user_service/pkg/helper"

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
	INSERT INTO "user" (
		id,
		first_name,
		last_name,
		phone_number,
		date_of_birth,
		password,
		created_at,
		updated_at
	) VALUES($1, $2, $3, $4, $5, $6, NOW(), NOW())
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

func (u *userRepo) GetByPK(ctx context.Context, req *user_service.UserPrimaryKey) (user *user_service.User, err error) {

	query := `
		SELECT 
			id,
			first_name,
			last_name,
			phone_number, 
			date_of_birth,
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
		Password:    password.String,
		CreatedAt:   created_at.String,
		UpdatedAt:   updated_at.String,
	}

	return user, nil
}

func (u *userRepo) GetList(ctx context.Context, req *user_service.GetListUserRequest) (resp *user_service.GetListUserResponse, err error) {
	resp = &user_service.GetListUserResponse{}

	var (
		query  string
		limit  = ""
		offset = " OFFSET 0 "
		params = make(map[string]interface{})
		filter = " WHERE TRUE "
		sort   = " ORDER BY created_at DESC"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			first_name,
			last_name,
			phone_number, 
			date_of_birth,
			password,
			TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "user"
	`
	if len(req.GetSearch()) > 0 {
		filter += " AND (first_name || ' ' || last_name) ILIKE '%' || '" + req.Search + "' || '%' "
	}
	if req.GetLimit() > 0 {
		limit = " LIMIT :limit"
		params["limit"] = req.Limit
	}
	if req.GetOffset() > 0 {
		offset = " OFFSET :offset"
		params["offset"] = req.Offset
	}
	query += filter + sort + offset + limit

	query, args := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
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

		err := rows.Scan(
			&resp.Count,
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
			return resp, err
		}

		resp.Users = append(resp.Users, &user_service.User{
			Id:          id.String,
			FirstName:   first_name.String,
			LastName:    last_name.String,
			PhoneNumber: phone_number.String,
			DateOfBirth: date_of_birth.String,
			Password:    password.String,
			CreatedAt:   created_at.String,
			UpdatedAt:   updated_at.String,
		})
	}

	return resp, nil
}

func (u *userRepo) Update(ctx context.Context, req *user_service.UpdateUser) (rowsAffected int64, err error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			"user"
		SET
			first_name = :first_name,
			last_name = :last_name,
			phone_number = :phone_number,
			date_of_birth = :date_of_birth,
			password = :password,
			updated_at = now()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"id":            req.GetId(),
		"first_name":    req.GetFirstName(),
		"last_name":     req.GetLastName(),
		"phone_number":  req.GetPhoneNumber(),
		"date_of_birth": req.GetDateOfBirth(),
		"password":      req.GetPassword(),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return
	}

	return result.RowsAffected(), nil
}

func (u *userRepo) UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (rowsAffected int64, err error) {

	var (
		set   = " SET "
		ind   = 0
		query string
	)

	if len(req.Fields) == 0 {
		err = errors.New("no updates provided")
		return
	}

	req.Fields["id"] = req.Id

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s ", key, key)
		if ind != len(req.Fields)-1 {
			set += ", "
		}
		ind++
	}

	query = `
		UPDATE
			"user"
	` + set + ` , updated_at = now()
		WHERE
			id = :id
	`

	query, args := helper.ReplaceQueryParams(query, req.Fields)

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return
	}

	return result.RowsAffected(), err
}

func (u *userRepo) Delete(ctx context.Context, req *user_service.UserPrimaryKey) error {

	query := `DELETE FROM "user" WHERE id = $1`

	_, err := u.db.Exec(ctx, query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
