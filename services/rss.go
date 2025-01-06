package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ogugu/models"
)

type RssService struct {
	db *sql.DB
}

func NewRssService(db *sql.DB) *RssService {
	return &RssService{
		db: db,
	}
}

func (r *RssService) DeleteRss(ctx context.Context, id string) error {
	spanctx, span := tracer.Start(ctx, "Delete RssFeed by ID")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := "DELETE FROM rss WHERE id = $1;"
	_, err := r.db.ExecContext(dbctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RssService) UpdateRss(ctx context.Context, id, field, value string) (models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "Update RssFeed Field")
	defer span.End()

	if field != "name" && field != "link" {
		return models.RssFeed{}, fmt.Errorf("Invalid field: %s. Only 'name' or 'link' are allowed", field)
	}

	var rss models.RssFeed
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`
		UPDATE rss
		SET %s = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, name, link, created_at, updated_at;
	`, field)

	row := r.db.QueryRowContext(dbctx, query, value, time.Now(), id)
	err := row.Scan(
		&rss.ID,
		&rss.Name,
		&rss.Link,
		&rss.CreatedAt,
		&rss.UpdatedAt,
	)
	if err != nil {
		return models.RssFeed{}, err
	}

	return rss, nil
}

func (r *RssService) GetAllRss(ctx context.Context) ([]models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "Get All RssFeeds")
	defer span.End()

	var allrss []models.RssFeed
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, name, link, created_at, updated_at FROM rss;`
	rows, err := r.db.QueryContext(dbctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var rss models.RssFeed
		err := rows.Scan(
			&rss.ID,
			&rss.Name,
			&rss.Link,
			&rss.CreatedAt,
			&rss.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		allrss = append(allrss, rss)
	}
	return allrss, nil
}

func (r *RssService) GetRss(ctx context.Context, field, value string) (models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "Get RssFeed by Field")
	defer span.End()

	if field != "name" && field != "link" {
		return models.RssFeed{}, fmt.Errorf("Invalid key: the %s field cannot be used as a key", field)
	}

	var rss models.RssFeed
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`SELECT id, name, link, created_at, updated_at FROM rss WHERE %s = $1;`, field)

	row := r.db.QueryRowContext(dbctx, query, value)
	err := row.Scan(
		&rss.ID,
		&rss.Name,
		&rss.Link,
		&rss.CreatedAt,
		&rss.UpdatedAt,
	)
	if err != nil {
		return models.RssFeed{}, err
	}

	return rss, nil
}

func (r *RssService) GetRssByID(ctx context.Context, id string) (models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "Get RssFeed by ID")
	defer span.End()

	var rss models.RssFeed
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, name, link, created_at, updated_at FROM rss WHERE id = $1;`

	row := r.db.QueryRowContext(dbctx, query, id)
	err := row.Scan(
		&rss.ID,
		&rss.Name,
		&rss.Link,
		&rss.CreatedAt,
		&rss.UpdatedAt,
	)
	if err != nil {
		return models.RssFeed{}, err
	}

	return rss, nil
}

func (r *RssService) CreateRSS(ctx context.Context, name, link, id string) (models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "Create RssFeed")
	defer span.End()

	var rss models.RssFeed

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		INSERT INTO rss (id, name, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, link, created_at, updated_at;
	`
	row := r.db.QueryRowContext(dbctx, query, id, name, link, time.Now(), time.Now())
	err := row.Scan(
		&rss.ID,
		&rss.Name,
		&rss.Link,
		&rss.CreatedAt,
		&rss.UpdatedAt,
	)
	if err != nil {
		return models.RssFeed{}, err
	}

	return rss, nil
}
