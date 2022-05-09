package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	return os.Getenv("MONGOURI")
}

func EnvHostName() string {
	val, ok := os.LookupEnv("HOSTNAME")
	if !ok {
		return "localhost"
	}
	return val
}

func EnvHostPort() string {
	val, ok := os.LookupEnv("HOSTNAME")
	if !ok {
		return "7000"
	}
	return val
}
