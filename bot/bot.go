package bot

import (
	"anon-chat-tg-bot/storage"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	b       *tgbotapi.BotAPI
	storage *storage.Store
	adminId int64
	chatIds map[int64]*storage.User
}

func New(token string, storage *storage.Store, adminId int64) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{
		b:       bot,
		storage: storage,
		adminId: adminId,
	}, nil
}

func (b *Bot) Run() error {
	//Get users from the database
	users, err := b.storage.GetUsers()
	if err != nil {
		panic(err)
	}

	//Create a map of users
	b.chatIds = make(map[int64]*storage.User)
	for _, u := range users {
		b.chatIds[u.ChatId] = &u
	}

	bot, err := tgbotapi.NewBotAPI(os.Args[1])
	if err != nil {
		return err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	if err := b.handleUpdates(updates); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		switch {
		case update.Message != nil && update.Message.IsCommand():
			if err := b.handleCommands(update); err != nil {
				return err
			}
		case update.Message != nil:
			b.sendToAll(*update.Message)
		}
	}
	return nil
}

func (b *Bot) handleCommands(update tgbotapi.Update) error {
	//Send the command to all users
	b.sendToAll(*update.Message)

	//Handle the command
	switch update.Message.Command() {
	case "start":
		//Check if the user is an admin
		root := false
		if update.Message.From.ID == b.adminId {
			root = true
		}

		//Create a new user
		newUser := storage.User{
			Id:       update.Message.From.ID,
			ChatId:   update.Message.Chat.ID,
			Username: update.Message.From.UserName,
			Root:     root,
		}

		//Add user to the map
		b.chatIds[update.Message.Chat.ID] = &newUser

		//Save the user to the database
		if err := b.storage.SaveUser(newUser); err != nil {
			return err
		}
	case "give_root":
		//Check if the user is an admin
		if update.Message.From.ID == b.adminId {
			//Get the chatId of the user to give root privellege
			for _, v := range b.chatIds {
				if v.Username == update.Message.CommandArguments() {
					v.Root = true
					if err := b.storage.SaveUser(*v); err != nil {
						return err
					}
				}
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are not an admin :/")
			m, err := b.b.Send(msg)
			if err != nil {
				return err
			}
			b.sendToAll(m)
		}
	case "remove_root":
		//Check if the user is an admin
		if update.Message.From.ID == b.adminId {
			//Get the chatId of the user to remove root privellege
			for _, v := range b.chatIds {
				if v.Username == update.Message.CommandArguments() && v.Root {
					v.Root = false
					if err := b.storage.SaveUser(*v); err != nil {
						return err
					}
				}
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are not an admin :/")
			m, err := b.b.Send(msg)
			if err != nil {
				return err
			}
			b.sendToAll(m)
		}
	case "count":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Number of users: %d", len(b.chatIds)))
		m, err := b.b.Send(msg)
		if err != nil {
			return err
		}
		b.sendToAll(m)
	}
	return nil
}

func (b *Bot) sendToAll(msg tgbotapi.Message) {
	sendersChatId := msg.Chat.ID
	for _, user := range b.chatIds {
		if user.ChatId != sendersChatId {
			if user.Root {
				forward := tgbotapi.NewForward(user.ChatId, sendersChatId, msg.MessageID)
				if _, err := b.b.Send(forward); err != nil {
					log.Printf("user: %s either blocked the bot or an error occured\n%v", user.Username, err)
				}
				continue
			}
			copy := tgbotapi.NewCopyMessage(user.ChatId, sendersChatId, msg.MessageID)
			if _, err := b.b.Send(copy); err != nil {
				log.Printf("user: %s either blocked the bot or an error occured\n%v", user.Username, err)
			}
		}
	}
}
