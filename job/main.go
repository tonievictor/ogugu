package job

import (
	"context"
	"database/sql"
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"

	"go.uber.org/zap"

	"ogugu/models"
	"ogugu/services/posts"
	"ogugu/services/rss"
)

func Start(db *sql.DB, logger *zap.Logger) {
	rssSrv := rss.New(db)

	feeds, err := rssSrv.Fetch(context.Background())
	if err != nil {
		logger.Error("could not get rss from db", zap.Error(err))
		return
	}

	for _, feed := range feeds {
		res, err := http.Get(feed.Link)
		if err != nil {
			logger.Error("an error occured while fetching rss data", zap.Error(err))
			continue
		}

		if !feed.Fetched {
			if err = populate(res, db, feed, logger); err != nil {
				continue
			}
			rssSrv.UpdateField(context.Background(), feed.ID, "fetched", true)
			continue
		}

		lm := res.Header.Get("Last-Modified")
		if lm == "" {
			continue
		}
		lastModified, err := time.Parse(time.RFC1123, lm)
		if err != nil {
			logger.Error("could not parse last modified time", zap.Error(err))
			continue
		}
		if lastModified.After(feed.LastModified) {
			if err = populate(res, db, feed, logger); err != nil {
				continue
			}
			rssSrv.UpdateField(context.Background(), feed.ID, "last_modified", lastModified)
		}
	}
}

func populate(res *http.Response, db *sql.DB, feed models.RssFeed, logger *zap.Logger) error {
	postSrv := posts.New(db)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("could not read response body", zap.Error(err))
		return err
	}

	var data models.RSSItems
	err = xml.Unmarshal(body, &data)
	if err != nil {
		logger.Error("could not unmarshal data "+feed.Link, zap.Error(err))
		return err
	}

	for _, value := range data.Channel.Items {
		_, err := postSrv.CreatePost(context.Background(), ulid.Make().String(), feed.ID, value)
		if err != nil {
			logger.Error("could not create a new post", zap.Error(err))
			continue
		}
	}

	return res.Body.Close()
}
