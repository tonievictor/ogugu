ALTER TABLE IF EXISTS subscriptions
ADD CONSTRAINT userid_rssid_unique_combo UNIQUE(user_id, rss_id);
