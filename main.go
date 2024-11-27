package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	chatIds := make(map[string]int64)
	//add Artem and Adres to the chat
	chatIds["fiodop"] = 667524408
	chatIds["dafraer"] = 2066065712
	chatIds["gonzalezthefast"] = 537161392
	bot, err := tgbotapi.NewBotAPI(os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			log.Printf("username: %s || userId: %d || chatId: %d\n", update.Message.From.UserName, update.Message.From.ID, update.Message.Chat.ID)
			sendersChatId := update.Message.Chat.ID
			sendersUserName := update.Message.From.UserName
			chatIds[sendersUserName] = sendersChatId
			for username, chatId := range chatIds {
				if chatId != sendersChatId {
					copy := tgbotapi.NewCopyMessage(chatId, sendersChatId, update.Message.MessageID)
					if _, err := bot.Send(copy); err != nil {
						log.Printf("user: %s either blocked the bot or error happened\n%v", username, err)
					}
				}
			}
		}
	}
}
