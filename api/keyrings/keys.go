package keyrings

import (
	"log"
	"strings"
	"sync"

	"github.com/TheAlpha16/typi/api/config"
)

var (
	apiKeys  []string
	keyIndex int = 0
	mu       sync.Mutex
)

func InitKeys() {
	keysString := config.YT_API_KEYS
	if keysString == "" {
		log.Fatal("No API keys found")
	}

	apiKeys = strings.Split(keysString, ",")
}

func GetKey() string {
	mu.Lock()
	defer mu.Unlock()
	
	currentKey := apiKeys[keyIndex]
	keyIndex = (keyIndex + 1) % len(apiKeys)
	return currentKey
}
