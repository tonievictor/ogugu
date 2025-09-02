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

var tracer = otel.Tracer("user service")

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (u *Service) GetUserByID(ctx context.Context, id string) (models.User, error) {
	spanctx, span := tracer.Start(ctx, "getuser by id")
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

func (u *Service) DeleteUserByID(ctx context.Context, id string) (int64, error) {
	spanctx, span := tracer.Start(ctx, "delete user")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1;`
	r, err := u.db.ExecContext(dbctx, query, id)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func (u *Service) CreateUser(ctx context.Context, id string, body models.CreateUserBody) (models.User, error) {
	spanctx, span := tracer.Start(ctx, "create user")
	defer span.End()

	var user models.User
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `INSERT into users (id, username, email, avatar, created_at, updated_at)
						VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, email, avatar, created_at, updated_at;
	`

	row := u.db.QueryRowContext(dbctx, query, id, body.Username, body.Email, body.Avatar, time.Now(), time.Now())
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

func (u *Service) UpdateUser(ctx context.Context, id string, field, value string) (models.User, error) {
	spanctx, span := tracer.Start(ctx, "update user")
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

func (u *Service) GetUser(ctx context.Context, field, value string) (models.User, error) {
	spanctx, span := tracer.Start(ctx, "fetch user")
	defer span.End()

	if field != "email" {
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

func (u *Service) GetUserAuth(ctx context.Context, email string) (string, string, error) {
	spanctx, span := tracer.Start(ctx, "get userid and password")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT users.id, auth.password FROM users INNER JOIN auth on users.id = auth.user_id WHERE users.email = $1;`
	var id string
	var password string
	row := u.db.QueryRowContext(dbctx, query, email)
	err := row.Scan(&id, &password)
	if err != nil {
		return "", "", err
	}

	return id, password, nil
}

func (u *Service) GetAllUsers(ctx context.Context) ([]models.User, error) {
	spanctx, span := tracer.Start(ctx, "fetch all users from db")
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
