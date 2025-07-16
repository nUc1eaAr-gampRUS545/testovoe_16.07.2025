package configs

import (
	"log"
	"os"
	"strconv"
)
import "github.com/joho/godotenv"

type Config struct {
	Port               string
	AllowedExtensions  []string
	MaxFilesInTask     int
	MaxConcurrentTasks int
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:               os.Getenv("PORT"),
		MaxFilesInTask:     getIntFromEnvFile("MAX_FILES"),
		MaxConcurrentTasks: getIntFromEnvFile("MAX_CONCURRENT_TASKS"),
		AllowedExtensions:  []string{".pdf", ".jpeg", ".jpg"},
	}
}

func getIntFromEnvFile(name string) int {
	numberStr := os.Getenv(name)
	if numberStr == "" {
		log.Fatalf(" %v not set in .env file", name)
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		log.Fatalf("Failed to convert %v to int: %v", name, err)
	}
	return number
}
