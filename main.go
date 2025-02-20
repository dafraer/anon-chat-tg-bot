package main

import (
	"anon-chat-tg-bot/bot"
	store "anon-chat-tg-bot/storage"
	"database/sql"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 4 {
		panic("Bot token, db uri and admin telegram user id expected as args")
	}
	//Create storage
	//default connection string for local postgres:
	//"user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	connStr := os.Args[2]
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		panic(err)
	}
	storage := store.New(db)
	if err := storage.Init(); err != nil {
		panic(err)
	}

	//Create the bot
	adminId, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Admin user id must be an integer: %v", err)
	}
	b, err := bot.New(os.Args[1], storage, int64(adminId))
	if err != nil {
		panic(err)
	}

	//Run the bot
	if err := b.Run(); err != nil {
		panic(err)
	}
}
