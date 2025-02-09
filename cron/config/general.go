package config

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var YT_API_KEYS = os.Getenv("YT_API_KEYS")
var TOPIC = "cricket"
var LAST_FETCH = time.Now().Add(-1 * time.Hour)
var FREQUENCY_IN_SECONDS = os.Getenv("FREQUENCY_IN_SECONDS")

func init() {
	if FREQUENCY_IN_SECONDS == "" {
		FREQUENCY_IN_SECONDS = "300"
	}
}
