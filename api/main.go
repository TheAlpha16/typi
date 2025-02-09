package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/TheAlpha16/typi/api/database"
	"github.com/TheAlpha16/typi/api/logs"
	"github.com/TheAlpha16/typi/api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

func main() {
	logs.InitLogger()

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

	// Setup access logs
	accessLogFile, err := os.OpenFile("./access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer accessLogFile.Close()
	aw := io.MultiWriter(os.Stdout, accessLogFile)
	loggerConfig := logger.Config{
		Output:     aw,
		TimeFormat: time.RFC3339,
	}

	app := fiber.New()
	app.Use(logger.New(loggerConfig))
	app.Use(recover.New())
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":4601"))
}
