package auth

import (
	"context"
	"database/sql"
	"time"

	"go.opentelemetry.io/otel"
)

var (
	dbtimeout = time.Second * 3
	tracer    = otel.Tracer("auth service")
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateAuth(ctx context.Context, id, password string) error {
	spanctx, span := tracer.Start(ctx, "Creating a new auth entry")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		INSERT INTO auth (user_id, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	_, err := r.db.ExecContext(dbctx, query, id, password, time.Now(), time.Now())
	return err
}

func (r *Repository) GetPasswordWithUserID(ctx context.Context, id string) (string, error) {
	spanctx, span := tracer.Start(ctx, "getting an auth entry")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	var password string
	query := `SELECT password FROM auth WHERE user_id = $1;`
	row := r.db.QueryRowContext(dbctx, query, id)
	err := row.Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}
