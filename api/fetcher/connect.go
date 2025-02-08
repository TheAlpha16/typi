package fetcher

import (
	"context"
	"log"

	"github.com/TheAlpha16/typi/api/keyrings"

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
		log.Println(err)
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
