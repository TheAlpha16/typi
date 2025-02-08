// // Sample Go code for user authorization

// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	// "golang.org/x/oauth2"
// 	// "google.golang.org/api/config/v1"
// 	"google.golang.org/api/option"
// 	"google.golang.org/api/youtube/v3"
// )

// func main() {
// 	ctx := context.Background()

// 	service, err := youtube.NewService(ctx, option.WithAPIKey("AIzaSyDOcTfdT1Z6QWrp8qqDjRl_c1sTabpfaiw"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	call, err := service.Search.List([]string{"snippet"}).
// 		Q("cats").
// 		MaxResults(25).
// 		Order("date").
// 		Type("video").
// 		Do()

// 	for _, item := range call.Items {
// 		jsonData, err := json.Marshal(item.Snippet)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(string(jsonData))
// 	}
// }

package main

import (
	"io"
	"log"
	"os"
	// "time"

	// "github.com/TheAlpha16/typi/api/database"
	"github.com/TheAlpha16/typi/api/logs"
	"github.com/TheAlpha16/typi/api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	logs.InitLogger()

	// for {
	// 	if err := database.Connect(); err != nil {
	// 		log.Println(err)
	// 		log.Println("sleep for 1 minute")
	// 		time.Sleep(time.Minute)
	// 		continue
	// 	}
	// 	break
	// }

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

	log.Fatal(app.Listen(":4601"))
}
