package main

import (
	"log"
	"os"
	"qask_telegram/internal/app/bot"
)

func main() {
	token := os.Getenv("TG_BOT_TOKEN")

	if token == "" {
		log.Fatal("Token not found")
	}

	if err := bot.Start(token); err != nil {
		log.Fatal(err)
	}
}
