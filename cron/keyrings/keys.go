/** 
Package keyrings provides functionality to manage and rotate API keys.

This package includes functions to initialize API keys from a configuration
and retrieve the current API key in a round-robin fashion.
*/

package keyrings

import (
	"strings"
	"sync"

	"github.com/TheAlpha16/typi/cron/config"

	"github.com/sirupsen/logrus"
)

var (
	apiKeys  []string
	keyIndex int = 0
	mu       sync.Mutex
)

/** 
InitKeys initializes the API keys from the configuration.
It retrieves the API keys from the config.YT_API_KEYS variable, which is expected
to be a comma-separated string of API keys. If no keys are found, the function
logs a fatal error and terminates the program.
*/

func InitKeys() {
	keysString := config.YT_API_KEYS
	if keysString == "" {
		logrus.Fatal("No API keys found")
	}

	apiKeys = strings.Split(keysString, ",")
}

/** 
GetKey returns the current API key and advances the key index in a round-robin
fashion. It ensures that the key retrieval and index update are thread-safe
by using a mutex lock.
*/

func GetKey() string {
	mu.Lock()
	defer mu.Unlock()

	currentKey := apiKeys[keyIndex]
	keyIndex = (keyIndex + 1) % len(apiKeys)
	return currentKey
}
