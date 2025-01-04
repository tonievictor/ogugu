package models

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("models.user")

type UserRepository interface {
	GetUser(ctx context.Context, field, value string) (User, error)
	DeleteUser(ctx context.Context, id string) error
	GetUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, username, email, id, avatar string) (User, error)
}
type User struct {
	ID        string    `json: "id"`
	Username  string    `json: "username" validate: "required"`
	Email     string    `json: "email" validate: "required:email"`
	Avatar    string    `json: "avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) DeleteUser(ctx context.Context, id string) error {
	_, span := tracer.Start(ctx, "MODELS delete user")
	defer span.End()

	dbctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
	defer cancel()

	query := ``
	_, err := db.ExecContext(dbctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CreateUser(ctx context.Context, username, email, id, avatar string) (User, error) {
	_, span := tracer.Start(ctx, "MODELS create user")
	defer span.End()

	var user User
	dbctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
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

func (u *User) GetUser(ctx context.Context, field, value string) (User, error) {
	_, span := tracer.Start(ctx, "MODELS getuser")
	defer span.End()
	var user User

	dbctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
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

func (u *User) GetUsers(ctx context.Context, field, value string) ([]User, error) {
	_, span := tracer.Start(ctx, "MODELS getusers")
	defer span.End()

	var users []User

	dbctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
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
