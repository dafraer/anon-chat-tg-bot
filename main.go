package main

import (
	store "adres-talk/storage"
	"database/sql"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	//Create storagee
	connStr := "user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	storage := store.New(db)
	if err := storage.Init(); err != nil {
		panic(err)
	}

	users, err := storage.GetUsers()
	if err != nil {
		panic(err)
	}

	chatIds := make(map[string]store.User)
	for _, u := range users {
		chatIds[u.Username] = u
	}

	bot, err := tgbotapi.NewBotAPI(os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message.IsCommand() {
			if update.Message.Text == "/start" {
				root := false
				if update.Message.From.UserName == "dafraer" {
					root = true
				}
				newUser := store.User{
					Id:       update.Message.From.ID,
					ChatId:   update.Message.Chat.ID,
					Username: update.Message.From.UserName,
					Root:     root,
				}
				if err := storage.SaveUser(newUser); err != nil {
					panic(err)
				}
				chatIds[update.Message.From.UserName] = newUser
			}
			continue
		}
		if update.Message != nil {
			log.Printf("username: %s || userId: %d || chatId: %d\n", update.Message.From.UserName, update.Message.From.ID, update.Message.Chat.ID)
			sendersChatId := update.Message.Chat.ID

			//check if the message is a command to give root privellege
			if chatIds[update.Message.From.UserName].Root && update.Message.Text[0] == '*' {
				nameOfNewRoot := update.Message.Text[1:]
				user, ok := chatIds[nameOfNewRoot]
				if ok {
					user.Root = true
					chatIds[nameOfNewRoot] = user
					if err := storage.SaveUser(user); err != nil {
						panic(err)
					}
				}
			}

			for _, user := range chatIds {
				if user.ChatId != sendersChatId {
					if user.Root {
						forward := tgbotapi.NewForward(user.ChatId, sendersChatId, update.Message.MessageID)
						if _, err := bot.Send(forward); err != nil {
							log.Printf("user: %s either blocked the bot or error happened\n%v", user.Username, err)
						}
						continue
					}
					copy := tgbotapi.NewCopyMessage(user.ChatId, sendersChatId, update.Message.MessageID)
					if _, err := bot.Send(copy); err != nil {
						log.Printf("user: %s either blocked the bot or error happened\n%v", user.Username, err)
					}
				}
			}
		}
	}
}
