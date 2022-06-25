package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host       string
	Port       string
	MongoDbUri string
}

func LoadConfigFromEnv() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return Config{
		Host:       getEnvOrDefault("APP_HOST", "localhost"),
		Port:       getEnvOrDefault("APP_PORT", "8080"),
		MongoDbUri: getEnvOrFail("MONGODB_URI", "You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable"),
	}
}

func getEnvOrDefault(variableName string, timeout string) string {
	var value string
	if value = os.Getenv(variableName); value == "" {
		return timeout
	}

	return value
}

func getEnvOrFail(variableName string, errorMessage string) string {
	var value string
	if value = os.Getenv(variableName); value == "" {
		log.Fatal(errorMessage)
		panic("Missing environment variable: " + variableName)
	}

	return value
}
