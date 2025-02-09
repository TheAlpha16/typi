package fetcher

import (
	"context"

	"github.com/TheAlpha16/typi/cron/keyrings"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	client *youtube.Service
)

func InitYTClient() *youtube.Service {
	apiKey := keyrings.GetKey()
	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logrus.WithError(err).Error("unable to create youtube service")
		return nil
	}

	return service
}

func GetYTClient() *youtube.Service {
	if client == nil {
		client = InitYTClient()
	}

	return client
}
