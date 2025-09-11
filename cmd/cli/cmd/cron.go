/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/oklog/ulid/v2"

	"ogugu/internal/models"
	"ogugu/internal/repository/posts"
	"ogugu/internal/repository/rss"
)

// cronCmd represents the cron command
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if err := job(dbConn); err == nil {
			fmt.Println("success!")
		}
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}

func job(db *sql.DB) error {
	rssSrv := rss.New(db)

	feeds, err := rssSrv.Fetch(context.Background())
	if err != nil {
		fmt.Println("could not get rss from db", err.Error())
		return err
	}

	for _, feed := range feeds {
		res, err := http.Get(feed.RSSLink)
		if err != nil {
			fmt.Println("an error occured while fetching rss data", err.Error())
			continue
		}

		if !feed.Fetched {
			if err = populate(res, db, feed); err != nil {
				continue
			}
			_, err := rssSrv.UpdateField(context.Background(), feed.ID, "fetched", true)
			if err != nil {
				fmt.Println("could not update fetched field ", err.Error())
			}
			continue
		}

		lm := res.Header.Get("Last-Modified")
		if lm == "" {
			continue
		}
		lastModified, err := time.Parse(time.RFC1123, lm)
		if err != nil {
			fmt.Println("could not parse last modified time", err.Error())
			continue
		}
		if lastModified.After(feed.LastModified) {
			if err = populate(res, db, feed); err != nil {
				continue
			}
			rssSrv.UpdateField(context.Background(), feed.ID, "last_modified", lastModified)
		}
	}
	return nil
}

func populate(res *http.Response, db *sql.DB, feed models.RssFeed) error {
	postSrv := posts.New(db)

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("could not read response body ", err.Error())
		return err
	}

	var data models.RSSItems
	err = xml.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("could not unmarshal data from "+feed.Link, err.Error())
		return err
	}

	for _, value := range data.Channel.Items {
		_, err := postSrv.CreatePost(context.Background(), ulid.Make().String(), feed.ID, value)
		if err != nil {
			fmt.Println("could not create a new post", err.Error())
			continue
		}
	}

	return nil
}
