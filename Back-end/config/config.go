package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	TgBotToken            string
	BasePath string
}

func MustLoad() Config {
	tgBotTokenToken := os.Getenv("TG_SKIN_DISEASE_CLASSIFIER_BOT")
	basePath := ""

	flag.Parse()

	if tgBotTokenToken == "" {
		log.Fatal("token is not specified")
	}


	return Config{
		TgBotToken:            tgBotTokenToken,
		BasePath: basePath,
	}
}
