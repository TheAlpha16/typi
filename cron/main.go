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
	"github.com/TheAlpha16/typi/cron/logs"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.InitLogger()
	keyrings.InitKeys()
	_ = fetcher.GetYTClient()

	for {
		if err := database.Connect(); err != nil {
			logrus.WithError(err).Error("failed to connect to database")
			logrus.Warn("sleep for 1 minute")
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

	_, err = c.AddFunc(fmt.Sprintf("@every %ss", config.FREQUENCY_IN_SECONDS), func() {
		logrus.Info("fetching new videos...")
		_ = fetcher.FetchVideosAsync()
	})

	if err != nil {
		logrus.Fatal("failed to schedule job: ", err)
	}

	c.Start()
	logrus.Info("fetcher started!")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logrus.Warn("shutting down fetcher...")

	ctx := c.Stop()
	select {
	case <-ctx.Done():
		logrus.Warn("cron scheduler stopped")
	case <-time.After(5 * time.Second):
		logrus.Error("timeout: force stopping scheduler")
	}

	logrus.Info("fetcher shut down")
}
