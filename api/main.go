package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/TheAlpha16/typi/api/config"
	"github.com/TheAlpha16/typi/api/database"
	"github.com/TheAlpha16/typi/api/fetcher"
	"github.com/TheAlpha16/typi/api/keyrings"
	"github.com/TheAlpha16/typi/api/logs"
	"github.com/TheAlpha16/typi/api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	logs.InitLogger()
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

	// Setup access logs
	accessLogFile, err := os.OpenFile("./access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer accessLogFile.Close()
	aw := io.MultiWriter(os.Stdout, accessLogFile)
	loggerConfig := logger.Config{
		Output: aw,
	}

	app := fiber.New()
	app.Use(logger.New(loggerConfig))
	app.Use(recover.New())
	router.SetupRoutes(app)

	// call fetch every 5 minutes
	go func() {
		for {
			log.Println("fetching new videos...")
			if err := fetcher.FetchVideosAsync(); err != nil {
				log.Println(err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	log.Fatal(app.Listen(":4601"))
}
