package main

import (
	"context"
	"dota_matches_bot/pkg/bot"
	"dota_matches_bot/pkg/repository/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	ctx := context.Background()
	if os.Getenv("ENV_RUNTIME") != "docker" {
		err := setLocalEnv()
		if err != nil {
			logrus.Fatalf("failed to load local env: %s", err.Error())
		}
	}

	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		logrus.Fatalf("failed to create botAPI: %v", err)
		return
	}

	logrus.Println("Telegram API connection established")

	rds, err := redis.NewRedisDb(ctx, redis.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err != nil {
		logrus.Fatalf("failed to open Redis connection, %v", err)
	}

	logrus.Println("Redis connection established")

	b := bot.NewBot(botAPI, rds)

	b.Handle(ctx)
}

func setLocalEnv() error {
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6381")
	if err := godotenv.Load(); err != nil {
		return err
	}

	return nil
}
