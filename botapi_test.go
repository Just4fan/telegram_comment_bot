package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"testing"
)

func TestBotAPI(t *testing.T) {
	token := "801920357:AAF0osBl_Znkci9D1dRw2NlgYIBTlffW62U"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
}

func TestStruct(t *testing.T) {
	data := []byte("add_comment eyJtZXNzYWdlX2lkIjoxNTUsImNoYXRfaWQiOjUxMDU0NjA2NX0=")
	log.Println(len(data))
}
