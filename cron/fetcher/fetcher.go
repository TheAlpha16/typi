package fetcher

import (
	"strings"
	"time"

	"github.com/TheAlpha16/typi/cron/config"
	"github.com/TheAlpha16/typi/cron/database"
	"github.com/TheAlpha16/typi/cron/models"
	"github.com/sirupsen/logrus"
)

func FetchVideosAsync() error {
	if client == nil {
		client = GetYTClient()
	}

	call := client.Search.List([]string{"snippet"}).
		Q(config.TOPIC).
		PublishedAfter(config.LAST_FETCH.Format(time.RFC3339)).
		Order("date").
		Type("video")

	response, err := call.Do()
	if err != nil {
		if strings.Contains(err.Error(), "quotaExceeded") {
			logrus.Warn("Quota exceeded, switching API key...")
			client = nil
			client = GetYTClient()
			return nil
		}
		logrus.WithError(err).Error("failed to fetch videos")
		return err
	}

	videos := make([]models.Video, 0)

	for _, item := range response.Items {
		video := models.Video{
			VID:         item.Id.VideoId,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			PublishedAt: item.Snippet.PublishedAt,
			Thumbnail:   item.Snippet.Thumbnails.Default.Url,
		}

		videos = append(videos, video)
	}

	if len(videos) == 0 {
		return nil
	}

	database.UploadVideos(&videos)

	lastPublishedAt, err := time.Parse(time.RFC3339, response.Items[0].Snippet.PublishedAt)
	if err != nil {
		logrus.WithError(err).Error("failed to parse last published at")
		return err
	}

	config.LAST_FETCH = lastPublishedAt
	database.UpdateLastFetch()

	return nil
}
