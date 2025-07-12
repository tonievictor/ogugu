ALTER TABLE IF EXISTS subscriptions
RENAME CONSTRAINT subscriptions_pkey TO users_rss_pkey;

ALTER TABLE IF EXISTS subscriptions
RENAME CONSTRAINT subscriptions_rss_id_fkey TO users_rss_rss_id_fkey;

ALTER TABLE IF EXISTS subscriptions
RENAME CONSTRAINT subscriptions_user_id_fkey TO users_rss_user_id_fkey;
