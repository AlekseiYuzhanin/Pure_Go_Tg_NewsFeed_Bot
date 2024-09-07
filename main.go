package main

import (
	"awesomeProject4/cleints/telegram"
	"flag"
	"log"
)

const (
	tgBotApi = "api.telegram.org"
)

func main() {

	//TODO: consumer start

	//TODO: init telegram bot(token)
	tgClient := telegram.New(tgBotApi, mustToken())
}

func mustToken() string {
	token := flag.String("token", "", "Token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("You must provide a token")
	}

	return *token
}
