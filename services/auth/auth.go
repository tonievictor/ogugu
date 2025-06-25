package auth

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel"
	"time"
)

var dbtimeout = time.Second * 3
var tracer = otel.Tracer("Auth Service")

type AuthService struct {
	db *sql.DB
}

func New(db *sql.DB) *AuthService {
	return &AuthService{db}
}

func (a *AuthService) CreateAuth(ctx context.Context, userId, password string) error {
	spanctx, span := tracer.Start(ctx, "Creating a new auth entry")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		INSERT INTO auth (user_id, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	_, err := a.db.ExecContext(dbctx, query, userId, password, time.Now(), time.Now())
	return err
}

func (a *AuthService) GetPassWordWithUserID(ctx context.Context, userId string) (string, error) {
	spanctx, span := tracer.Start(ctx, "getting an auth entry")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	var password string
	query := `SELECT password FROM auth WHERE user_id = $1;`
	row := a.db.QueryRowContext(dbctx, query, userId)
	err := row.Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}
