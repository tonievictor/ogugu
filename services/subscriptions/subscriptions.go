package subscriptions

import (
	"context"
	"database/sql"
	"time"

	"go.opentelemetry.io/otel"
	"ogugu/models"
)

const dbtimeout = time.Second * 3

var tracer = otel.Tracer("rss service")

type SubscriptionService struct {
	db *sql.DB
}

func New(db *sql.DB) *SubscriptionService {
	return &SubscriptionService{
		db: db,
	}
}

func (ss *SubscriptionService) CreateSub(ctx context.Context, id, user_id, rss_id string) (models.Subscription, error) {
	spanctx, span := tracer.Start(ctx, "create a new subscription")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		INSERT INTO subscriptions (id, user_id, rss_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, created_at, updated_at;
	`

	var sub models.Subscription
	row := ss.db.QueryRowContext(dbctx, query, id, user_id, rss_id, time.Now(), time.Now())
	err := row.Scan(&sub.ID, &sub.UserID, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return models.Subscription{}, err
	}

	return sub, nil
}

func (ss *SubscriptionService) GetSubByID(ctx context.Context, id string) (models.Subscription, error) {
	spanctx, span := tracer.Start(ctx, "get subscription by id")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		SELECT sub.id, sub.user_id, sub.created_at, sub.updated_at,
		rss.id, rss.title, rss.link, rss.created_at, rss.updated_at 
		FROM subscriptions sub
		INNER JOIN rss ON rss.id = sub.rss_id
		WHERE sub.id = $1;
	`

	var sub models.Subscription
	row := ss.db.QueryRowContext(dbctx, query, id)
	err := row.Scan(
		&sub.ID,
		&sub.UserID,
		&sub.CreatedAt,
		&sub.UpdatedAt,
		&sub.RSS.ID,
		&sub.RSS.Title,
		&sub.RSS.Link,
		&sub.RSS.CreatedAt,
		&sub.RSS.UpdatedAt,
	)
	if err != nil {
		return models.Subscription{}, err
	}

	return sub, nil
}

func (ss *SubscriptionService) GetSubs(ctx context.Context) ([]models.Subscription, error) {
	spanctx, span := tracer.Start(ctx, "get all subscriptions")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		SELECT sub.id, sub.user_id, sub.created_at, sub.updated_at,
		rss.id, rss.title, rss.link, rss.created_at, rss.updated_at 
		FROM subscriptions sub
		INNER JOIN rss ON rss.id = sub.rss_id;
	`
	rows, err := ss.db.QueryContext(dbctx, query)
	if err != nil {
		return nil, err
	}

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&sub.RSS.ID,
			&sub.RSS.Title,
			&sub.RSS.Link,
			&sub.RSS.CreatedAt,
			&sub.RSS.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (ss *SubscriptionService) GetSubsByUserID(ctx context.Context, user_id string) ([]models.Subscription, error) {
	spanctx, span := tracer.Start(ctx, "get all subscriptions by user id")
	defer span.End()

	dbctx, cancel := context.WithTimeout(spanctx, dbtimeout)
	defer cancel()

	query := `
		SELECT sub.id, sub.user_id, sub.created_at, sub.updated_at,
		rss.id, rss.title, rss.link, rss.created_at, rss.updated_at 
		FROM subscriptions sub
		INNER JOIN rss ON rss.id = sub.rss_id
		WHERE sub.user_id = $1;
	`
	rows, err := ss.db.QueryContext(dbctx, query, user_id)
	if err != nil {
		return nil, err
	}

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.CreatedAt,
			&sub.UpdatedAt,
			&sub.RSS.ID,
			&sub.RSS.Title,
			&sub.RSS.Link,
			&sub.RSS.CreatedAt,
			&sub.RSS.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}
