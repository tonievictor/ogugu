package models

import (
	"context"
	"fmt"
	"time"
)

type RssRepository interface {
	CreateRSS(ctx context.Context, name, link, id string) (Rss, error)
	GetRssByID(ctx context.Context, id string) (Rss, error)
	GetRss(ctx context.Context, field, value string) (Rss, error)
	GetAllRss(ctx context.Context) ([]Rss, error)
	UpdateRss(ctx context.Context, id, field, value string) (Rss, error)
	DeleteRss(ctx context.Context, id string) error
}

type Rss struct {
	ID        string    `json: "id" validate: "required"`
	Name      string    `json: "name" validate: "required"`
	Link      string    `json: "link" validate: "required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r Rss) DeleteRss(ctx context.Context, id string) error {
	spanctx, span := tracer.Start(ctx, "Delete Rss")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := "DELETE FROM rss WHERE id = $1;"
	_, err := db.ExecContext(dbctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r Rss) UpdateRss(ctx context.Context, id, field, value string) (Rss, error) {
	spanctx, span := tracer.Start(ctx, "Update Rss")
	defer span.End()

	if field != "name" && field != "link" {
		return Rss{}, fmt.Errorf("Cannot update the %s field", field)
	}

	var rss Rss
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`
		UPDATE rss
		SET %s = $1,
		WHERE id = $2,
		RETURNING id, name, link, created_at, updated_at;
	`, field)

	row := db.QueryRowContext(dbctx, query, value, id)
	err := row.Scan(
		&rss.ID,
		&rss.Name,
		&rss.Link,
		&rss.CreatedAt,
		&rss.CreatedAt,
	)
	if err != nil {
		return Rss{}, err
	}

	return rss, nil
}

func (r Rss) GetAllRss(ctx context.Context) ([]Rss, error) {
	spanctx, span := tracer.Start(ctx, "Get all Rss")
	defer span.End()

	var allrss []Rss
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, name, link, created_at, updated_at FROM rss;`
	rows, err := db.QueryContext(dbctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rss Rss
		err := rows.Scan(
			&rss.ID,
			&rss.Name,
			&rss.Link,
			&rss.CreatedAt,
			&rss.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		allrss = append(allrss, rss)
	}
	return allrss, nil
}

func (r Rss) GetRss(ctx context.Context, field, value string) (Rss, error) {
	spanctx, span := tracer.Start(ctx, "Get Rss")
	defer span.End()

	if field != "name" && field != "link" {
		return Rss{}, fmt.Errorf("Cannot update the %s field", field)
	}

	var rss Rss
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`SELECT id, name, link, created_at, updated_at FROM rss WHERE %s = $1;`, field)

	row := db.QueryRowContext(dbctx, query, value)
	err := row.Scan(
		&rss.ID,
		&rss.Name,
		&rss.Link,
		&rss.CreatedAt,
		&rss.CreatedAt,
	)
	if err != nil {
		return Rss{}, err
	}

	return rss, nil
}

func (r Rss) GetRssByID(ctx context.Context, id string) (Rss, error) {
	spanctx, span := tracer.Start(ctx, "Get Rss by ID")
	defer span.End()

	var rss Rss
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, name, link, created_at, updated_at FROM rss WHERE id = $1;`

	row := db.QueryRowContext(dbctx, query, id)
	err := row.Scan(
		&rss.ID,
		&rss.Name,
		&rss.Link,
		&rss.CreatedAt,
		&rss.CreatedAt,
	)
	if err != nil {
		return Rss{}, err
	}

	return rss, nil
}

func (r *Rss) CreateRSS(ctx context.Context, name, link, id string) (Rss, error) {
	spanctx, span := tracer.Start(ctx, "Create rss")
	defer span.End()

	var rss Rss

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		INSERT INTO rss (id, name, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, link, created_at, updated_at;
	`
	row := db.QueryRowContext(dbctx, query, id, name, link, time.Now(), time.Now())
	err := row.Scan(
		&rss.ID,
		&rss.Name,
		&rss.Link,
		&rss.CreatedAt,
		&rss.CreatedAt,
	)
	if err != nil {
		return Rss{}, err
	}

	return rss, nil
}
