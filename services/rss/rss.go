package rss

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"ogugu/models"
)

const dbtimeout = time.Second * 3

var tracer = otel.Tracer("rss service")

type RssService struct {
	db *sql.DB
}

func New(db *sql.DB) *RssService {
	return &RssService{
		db: db,
	}
}

func (r *RssService) DeleteByID(ctx context.Context, id string) (int64, error) {
	spanctx, span := tracer.Start(ctx, "delete rss feed by id")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := "DELETE FROM rss WHERE id = $1;"
	res, err := r.db.ExecContext(dbctx, query, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *RssService) UpdateField(ctx context.Context, id, field, value any) (models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "update rss feed")
	defer span.End()

	var rss models.RssFeed
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := fmt.Sprintf(`
		UPDATE rss
		SET %s = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, title, link, description, fetched, last_modified, created_at, updated_at;
	`, field)

	row := r.db.QueryRowContext(dbctx, query, value, time.Now(), id)
	err := row.Scan(
		&rss.ID,
		&rss.Title,
		&rss.Link,
		&rss.Description,
		&rss.Fetched,
		&rss.LastModified,
		&rss.CreatedAt,
		&rss.UpdatedAt,
	)
	if err != nil {
		return models.RssFeed{}, err
	}

	return rss, nil
}

func (r *RssService) Fetch(ctx context.Context) ([]models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "fetch all rss feeds")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, title, link, description, fetched, last_modified, created_at, updated_at FROM rss;`
	rows, err := r.db.QueryContext(dbctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allrss []models.RssFeed
	for rows.Next() {
		var rss models.RssFeed
		err := rows.Scan(
			&rss.ID,
			&rss.Title,
			&rss.Link,
			&rss.Description,
			&rss.Fetched,
			&rss.LastModified,
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

func (r *RssService) FindByID(ctx context.Context, id string) (models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "fetch rss feed by id")
	defer span.End()

	var rss models.RssFeed
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, title, link, description, fetched, last_modified, created_at, updated_at FROM rss WHERE id = $1;`

	row := r.db.QueryRowContext(dbctx, query, id)
	err := row.Scan(
		&rss.ID,
		&rss.Title,
		&rss.Link,
		&rss.Description,
		&rss.Fetched,
		&rss.LastModified,
		&rss.CreatedAt,
		&rss.UpdatedAt,
	)
	if err != nil {
		return models.RssFeed{}, err
	}

	return rss, nil
}

func (r *RssService) FindByLink(ctx context.Context, link string) (models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "fetch rss feed by link")
	defer span.End()

	var rss models.RssFeed
	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `SELECT id, title, link, description, fetched, last_modified, created_at, updated_at FROM rss WHERE link = $1;`

	row := r.db.QueryRowContext(dbctx, query, link)
	err := row.Scan(
		&rss.ID,
		&rss.Title,
		&rss.Link,
		&rss.Description,
		&rss.Fetched,
		&rss.LastModified,
		&rss.CreatedAt,
		&rss.UpdatedAt,
	)
	if err != nil {
		return models.RssFeed{}, err
	}

	return rss, nil
}

func (r *RssService) Create(ctx context.Context, id, link string, body models.RSSMeta) (models.RssFeed, error) {
	spanctx, span := tracer.Start(ctx, "insert rss feed")
	defer span.End()

	var rss models.RssFeed

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		INSERT INTO rss (id, title, link, description, last_modified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, link, description, fetched, last_modified, created_at, updated_at;
	`
	row := r.db.QueryRowContext(dbctx, query, id, body.Channel.Title, link, body.Channel.Description, body.Channel.LastModified, time.Now(), time.Now())
	err := row.Scan(
		&rss.ID,
		&rss.Title,
		&rss.Link,
		&rss.Description,
		&rss.Fetched,
		&rss.LastModified,
		&rss.CreatedAt,
		&rss.UpdatedAt,
	)
	if err != nil {
		return models.RssFeed{}, err
	}

	return rss, nil
}
