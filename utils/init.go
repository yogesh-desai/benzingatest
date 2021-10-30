package utils

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// Init initialises basic data structures required.
func Init() {

	dir, _ := os.Getwd()
	err := godotenv.Load(dir + "/files/config/main.env")
	if err != nil {
		log.Fatal("[ERROR][Init] Unable to load .env fie. ", err)
	}

	log.Println("[INFO][Init] Environment loaded successfully.")

	Cfg = GetConfig()

	log.Printf("[INFO][Init] Config populated. %+v\n", Cfg)

	// init cache according to batch size
	Cache, err = initCache(Cfg.BatchSize)
	if err != nil {
		log.Fatal("[ERROR][Init] Unable to init cache.", err.Error())
	}
}

type Count struct {
	Mu  sync.RWMutex
	Key int
}

var KeyCount *Count = &Count{}

// GetCount returns the count
func GetCount() *Count {
	if KeyCount == nil {
		KeyCount = &Count{}
	}

	return KeyCount
}

// IncrementCacheCount increments the counter used for cache management
func IncrementCacheCount() {
	KeyCount.Mu.Lock()
	KeyCount.Key = KeyCount.Key + 1
	KeyCount.Mu.Unlock()
}

// ResetCacheCount resets the cache
func ResetCacheCount() {
	KeyCount.Mu.Lock()
	KeyCount.Key = 0
	KeyCount.Mu.Unlock()
}

// Config struct
type Config struct {
	BatchSize     int
	BatchInterval int
	Endpoint      string
}

var Cfg *Config

// GetConfig returns the config
func GetConfig() *Config {

	var config *Config = &Config{}

	env, ok := os.LookupEnv("batchSize")
	if !ok {
		log.Println("[INFO][GetConfig] Environment variable - BatchSize not set, using defaults")
		config.BatchSize = 1
	} else {

		size, err := strconv.Atoi(env)
		if err != nil {
			log.Println("[ERROR][GetConfig] Error in string conversion. ", err)
		}

		config.BatchSize = size
	}

	env, ok = os.LookupEnv("batchInterval")
	if !ok {
		log.Println("[INFO][GetConfig] Environment variable - batchInterval not set, using defaults")
		config.BatchInterval = 2
	} else {
		interval, err := strconv.Atoi(env)
		if err != nil {
			log.Println("[ERROR][GetConfig] Error in string conversion. ", err)
		}

		config.BatchInterval = interval
	}

	env, ok = os.LookupEnv("postEndpoint")
	if !ok {
		log.Println("[INFO][GetConfig] Environment variable - postEndpoint not set, using defaults")
		config.Endpoint = "https://enjk446ttn3h.x.pipedream.net/" // This I have created as public request bin.
	} else {

		if env != "" {
			config.Endpoint = env
		}
	}

	return config
}
