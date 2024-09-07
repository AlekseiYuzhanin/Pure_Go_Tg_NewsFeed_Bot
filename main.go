package main

import (
	tgClient "awesomeProject4/cleints/telegram"
	event_consumer "awesomeProject4/consumer/event-consumer"
	"awesomeProject4/events/telegram"
	"awesomeProject4/storage/files"
	"flag"
	"log"
)

// 7518956179:AAEgEMMtXooApEtke6ojZmpu5forLOBAA34
const (
	tgBotApi    = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotApi, mustToken()),
		files.New(storagePath),
	)
	log.Println("Starting telegram bot")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

func mustToken() string {
	token := flag.String("token", "", "Token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("You must provide a token")
	}

	return *token
}
