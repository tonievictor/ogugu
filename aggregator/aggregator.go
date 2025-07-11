package aggregator

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
	"ogugu/models"
)

type rssMeta struct {
	Channel struct {
		LastBuildDate string `xml:"lastBuildDate"`
		Title         string `xml:"title" validate:"required"`
		Description   string `xml:"description" validate:"required"`
		Link          string `xml:"link" validate:"required,url"`
		Items         []models.RssFeed
	} `xml:"channel"`
}

func RssToPosts(url string, log *zap.Logger) error {
	client := http.Client{
		Timeout: time.Second * 60,
	}

	res, err := client.Get(url)
	if err != nil {
		log.Error("cannot get response body", zap.Error(err))
		return err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("cannot read response body", zap.Error(err))
		return err
	}

	var body rssMeta
	err = xml.Unmarshal(data, &body)
	if err != nil {
		log.Error("cannot unmarshal rss feed", zap.Error(err))
		return err
	}

	return nil
}
