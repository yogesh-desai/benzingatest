package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Init initialises basic data structures required.
func Init() {

	dir, _ := os.Getwd()
	err := godotenv.Load(dir + "/files/config/main.env")
	if err != nil {
		log.Println("[Init] Unable to load .env fie. ", err)
	}

	log.Println("[Init] Environment loaded successfully.")

	Cfg = GetConfig()

	log.Printf("[Init] Config populated. %+v\n", Cfg)
}

var Cfg *Config

type Config struct {
	BatchSize     int
	BatchInterval int
	Endpoint      string
}

func GetConfig() *Config {

	var config *Config

	env, ok := os.LookupEnv("batchSize")
	if !ok {
		log.Println("[GetConfig] Environment variable - BatchSize not set, using defaults")
		config.BatchSize = 1
	} else {

		size, err := strconv.Atoi(env)
		if err != nil {
			log.Println("[GetConfig] Error in string conversion. ", err)
		}

		config.BatchSize = size
	}

	env, ok = os.LookupEnv("batchInterval")
	if !ok {
		log.Println("[GetConfig] Environment variable - batchInterval not set, using defaults")
		config.BatchInterval = 2
	} else {
		interval, err := strconv.Atoi(env)
		if err != nil {
			log.Println("[GetConfig] Error in string conversion. ", err)
		}

		config.BatchSize = interval
	}

	env, ok = os.LookupEnv("postEndpoint")
	if !ok {
		log.Println("[GetConfig] Environment variable - postEndpoint not set, using defaults")
		config.Endpoint = "/log"
	} else {
		config.Endpoint = env
	}

	return config
}
