package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheAlpha16/typi/cron/config"
	"github.com/TheAlpha16/typi/cron/database"
	"github.com/TheAlpha16/typi/cron/fetcher"
	"github.com/TheAlpha16/typi/cron/keyrings"

	"github.com/robfig/cron/v3"
)

func main() {
	keyrings.InitKeys()
	_ = fetcher.GetYTClient()

	for {
		if err := database.Connect(); err != nil {
			log.Println(err)
			log.Println("sleep for 1 minute")
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	last_fetch, err := database.GetLastFetch()
	if err != nil {
		log.Fatal(err)
	}
	config.LAST_FETCH = last_fetch

	c := cron.New()

	_, err = c.AddFunc(fmt.Sprintf("@every %ds", config.FREQUENCY_IN_SECONDS), func() {
		log.Println("fetching new videos...")
		if err := fetcher.FetchVideosAsync(); err != nil {
			log.Println(err)
		}
	})

	if err != nil {
		log.Fatal("failed to schedule job: ", err)
	}

	c.Start()
	log.Println("fetcher started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("shutting down fetcher...")

	ctx := c.Stop()
	select {
	case <-ctx.Done():
		log.Println("cron scheduler stopped")
	case <-time.After(5 * time.Second):
		log.Println("timeout: force stopping scheduler")
	}

	log.Println("fetcher shut down")
}
