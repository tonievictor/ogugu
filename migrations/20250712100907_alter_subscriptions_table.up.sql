ALTER TABLE IF EXISTS subscriptions
RENAME CONSTRAINT users_rss_pkey TO subscriptions_pkey;

ALTER TABLE IF EXISTS subscriptions
RENAME CONSTRAINT users_rss_rss_id_fkey TO subscriptions_rss_id_fkey;

ALTER TABLE IF EXISTS subscriptions
RENAME CONSTRAINT users_rss_user_id_fkey TO subscriptions_user_id_fkey;
