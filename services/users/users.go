package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"

	"ogugu/models"
)

const dbtimeout = time.Second * 3

var tracer = otel.Tracer("User Service")

type UserService struct {
	db *sql.DB
}

func New(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (u *UserService) GetUserByID(ctx context.Context, id string) (models.User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS getuser by id")
	defer span.End()

	var user models.User

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, username, email, avatar, created_at, updated_at FROM users WHERE id = $1;`
	row := u.db.QueryRowContext(dbctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UserService) DeleteUser(ctx context.Context, id string) error {
	spanctx, span := tracer.Start(ctx, "MODELS delete user")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1;`
	_, err := u.db.ExecContext(dbctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) CreateUser(ctx context.Context, username, email, id, avatar string) (models.User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS create user")
	defer span.End()

	var user models.User
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `INSERT into users (id, username, email, avatar, created_at, updated_at)
						VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, email, avatar, created_at, updated_at;
	`

	row := u.db.QueryRowContext(dbctx, query, id, username, email, avatar, time.Now(), time.Now())
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UserService) UpdateUser(ctx context.Context, id string, field, value string) (models.User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS update user")
	defer span.End()

	if field != "email" && field != "username" && field != "avatar" {
		return models.User{}, fmt.Errorf("field %s cannot be updated", field)
	}

	var user models.User

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`
		UPDATE users
		SET %s = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, username, email, avatar, created_at, updated_at;`, field)

	row := u.db.QueryRowContext(dbctx, query, value, time.Now(), id)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UserService) GetUser(ctx context.Context, field, value string) (models.User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS getuser")
	defer span.End()

	if field != "email" && field != "username" && field != "avatar" {
		return models.User{}, fmt.Errorf("Cannot use the %s field as a key", field)
	}

	var user models.User

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`SELECT id, username, email, avatar, created_at, updated_at FROM users WHERE %s = $1;`, field)
	row := u.db.QueryRowContext(dbctx, query, value)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	spanctx, span := tracer.Start(ctx, "MODELS getusers")
	defer span.End()

	var users []models.User

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, username, email, avatar, created_at, updated_at FROM users;`
	rows, err := u.db.QueryContext(dbctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
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
