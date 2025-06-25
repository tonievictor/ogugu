package posts

import (
	"context"
	"database/sql"
	"ogugu/models"
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("RssFeed Service")

const dbtimeout = time.Second * 3

type PostService struct {
	db *sql.DB
}

func New(db *sql.DB) *PostService {
	return &PostService{db: db}
}

func (ps *PostService) CreatePost(ctx context.Context, id string, rss_id string, p models.CreatePost) (models.Post, error) {
	spanctx, span := tracer.Start(ctx, "creating a new post")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		INSERT INTO posts (id, rss_id, title, description, link, pubdate, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, link, pubdate, created_at, updated_at;
	`
	row := ps.db.QueryRowContext(dbctx, query, id, rss_id, p.Title, p.Description, p.Link, p.PubDate, time.Now(), time.Now())

	var post models.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Description,
		&post.Link,
		&post.PubDate,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}
