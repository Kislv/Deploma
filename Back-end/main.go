package main

import (
	tgClient "read-adviser-bot/clients/telegram"
	"read-adviser-bot/config"
	"read-adviser-bot/utils/log"
	"read-adviser-bot/consumer/event-consumer"
	"read-adviser-bot/events/telegram"
	json_config "read-adviser-bot/utils/config"
	database "read-adviser-bot/storage/PostgreSQL"
)

const (
	tgBotHost   = "api.telegram.org"
	tgBotFileStorage   = tgBotHost + "/file"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	err := json_config.DevConfigStore.FromJson()
	if err != nil {
		log.Error(err)
	}

	err = json_config.ProdConfigStore.FromJson()
	if err != nil {
		log.Error(err)
	}


	cfg := config.MustLoad()


	db := database.InitDatabase()
	db.Connect()
	defer db.Disconnect()

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TgBotToken),
		db,
	)

	println("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		println("service is stopped")
	}
}
