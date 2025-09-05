package posts

import (
	"context"
	"database/sql"
	"time"

	"go.opentelemetry.io/otel"
	"ogugu/models"
)

var tracer = otel.Tracer("posts service")

const dbtimeout = time.Second * 3

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePost(
	ctx context.Context, id string, rss_id string, p models.CreatePost,
) (models.Post, error) {
	spanctx, span := tracer.Start(ctx, "creating a new post")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		INSERT INTO posts (id, rss_id, title, description, link, pubdate, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, link, pubdate, created_at, updated_at;
	`
	row := r.db.QueryRowContext(
		dbctx, query, id, rss_id, p.Title, p.Description, p.Link, p.PubDate, time.Now(), time.Now(),
	)

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

func (r *Repository) GetByID(ctx context.Context, id string) (models.Post, error) {
	spanctx, span := tracer.Start(ctx, "get a post by id")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		SELECT id, title, description, link, pubdate, created_at, updated_at 
		FROM posts WHERE id = $1;
	`
	row := r.db.QueryRowContext(dbctx, query, id)
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

func (r *Repository) Fetch(ctx context.Context) ([]models.Post, error) {
	spanctx, span := tracer.Start(ctx, "fetch all posts")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		SELECT id, title, description, link, pubdate, created_at, updated_at FROM posts;
	`
	rows, err := r.db.QueryContext(dbctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Description,
			&post.Link,
			&post.PubDate,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (ps *Repository) DeletePost(ctx context.Context, id string) (int64, error) {
	spanctx, span := tracer.Start(ctx, "delete post by id")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `DELETE FROM posts WHERE id = $1;`
	r, err := ps.db.ExecContext(dbctx, query, id)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
