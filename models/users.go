package models

import (
	"context"
	"fmt"
	"time"
)

type UserRepository interface {
	GetUser(field, value string) (*User, error)
	GetUsers() ([]*User, error)
}

type User struct {
	ID        string    `json: "id"`
	Username  string    `json: "username" validate: "required"`
	Email     string    `json: "email" validate: "required:email"`
	Avatar    string    `json: "avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetUser(field, value string) (User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`SELECT id, username, email, avatar, created_at, updated_at FROM users WHERE %s = $1`, field)

	row := db.QueryRowContext(ctx, query, value)

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
