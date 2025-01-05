package models

import (
	"context"
	"fmt"
	"time"
)

type UserRepository interface {
	GetUser(ctx context.Context, field, value string) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, field, value string) (User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, username, email, id, avatar string) (User, error)
}

type User struct {
	ID        string    `json: "id" validate: "required"`
	Username  string    `json: "username" validate: "required"`
	Email     string    `json: "email" validate: "required:email"`
	Avatar    string    `json: "avatar" validate: "required"`
	CreatedAt time.Time `json:"created_at" validate: "required"`
	UpdatedAt time.Time `json:"updated_at" validate: "required"`
}

func (u User) GetUserByID(ctx context.Context, id string) (User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS getuser by id")
	defer span.End()

	var user User

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, username, email, avatar, created_at, updated_at FROM users WHERE id = $1;`
	row := db.QueryRowContext(dbctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, nil
	}

	return user, nil
}

func (u User) DeleteUser(ctx context.Context, id string) error {
	spanctx, span := tracer.Start(ctx, "MODELS delete user")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1;`
	_, err := db.ExecContext(dbctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u User) CreateUser(ctx context.Context, username, email, id, avatar string) (User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS create user")
	defer span.End()

	var user User
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `INSERT into users (id, username, email, avatar, created_at, updated_at)
						VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, email, avatar, created_at, updated_at;
	`

	row := db.QueryRowContext(dbctx, query, id, username, email, avatar, time.Now(), time.Now())
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (u User) UpdateUser(ctx context.Context, id string, field, value string) (User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS update user")
	defer span.End()

	if field == "email" && field != "username" && field != "avatar" {
		return User{}, fmt.Errorf("Cannot update the %s field", field)
	}

	var user User

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`
		UPDATE users 
		SET %s = $1
		SET updated_at = $2
		WHERE id = $2;`, field)
	row := db.QueryRowContext(dbctx, query, value, id, time.Now())

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, nil
	}

	return user, nil
}

func (u User) GetUser(ctx context.Context, field, value string) (User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS getuser")
	defer span.End()

	if field == "email" && field != "username" && field != "avatar" {
		return User{}, fmt.Errorf("Cannot update the %s field", field)
	}

	var user User

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`SELECT id, username, email, avatar, created_at, updated_at FROM users WHERE %s = $1;`, field)
	row := db.QueryRowContext(dbctx, query, value)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, nil
	}

	return user, nil
}

func (u User) GetAllUsers(ctx context.Context) ([]User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS getusers")
	defer span.End()

	var users []User

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, username, email, avatar, created_at, updated_at FROM users;`
	rows, err := db.QueryContext(dbctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Avatar,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
