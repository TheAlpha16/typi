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

/*
main is the entry point of the application. It initializes the logger, keyrings, and YouTube client.
It then attempts to connect to the database, retrying every minute if the connection fails.
Once connected, it retrieves the last fetch time from the database and sets it in the configuration.
A new cron scheduler is created to fetch new videos at intervals specified by the configuration.
The application listens for interrupt signals to gracefully shut down the cron scheduler and the fetcher.
*/

func main() {
	logs.InitLogger()
	keyrings.InitKeys()
	_ = fetcher.GetYTClient()
	
	logrus.Info("connecting to database...")
	for {
		if err := database.Connect(); err != nil {
			logrus.WithError(err).Error("failed to connect to database")
			logrus.Warn("sleep for 1 minute")
			time.Sleep(time.Minute)
			continue
		}
		break
	}
	logrus.Info("connected to database")

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
