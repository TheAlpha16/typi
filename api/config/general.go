package config

import (
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var YT_API_KEYS = os.Getenv("YT_API_KEYS")
var TOPIC = "cricket"
var LAST_FETCH = time.Now().Add(-1 * time.Hour)
var DEFAULT_PER_PAGE = 10
var MAX_PER_PAGE = 100
