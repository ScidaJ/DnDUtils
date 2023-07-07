package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvAPIURL() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("API_ENDPOINT")
}

func EnvBotToken() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("BOT_TOKEN")
}

func EnvGuildId() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("GUILD_ID")
}
