package main

import (
	"flag"
	"log"
	"os"
	"qask_telegram/internal/app/bot"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "./configs/qask_telegram.conf", "path to config file")
}

func main() {
	flag.Parse()

	token := os.Getenv("TG_BOT_TOKEN")

	config := bot.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if token == "" {
		log.Fatal("Token not found")
	}

	if err := bot.Start(token, config); err != nil {
		log.Fatal(err)
	}
}
