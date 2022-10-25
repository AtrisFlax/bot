package main

import (
	"bot/internal/service/product"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	u := tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 60,
	}

	updates := bot.GetUpdatesChan(u)

	productService := product.NewService()

	for update := range updates {
		if update.Message != nil { // If we got a message
			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "help":
				helpCommand(bot, update.Message)
			case "list":
				listCommand(bot, update.Message, productService)
			default:
				defaultBehavior(bot, update.Message)
			}
		}
	}
}

func listCommand(bot *tgbotapi.BotAPI, inputMessage *tgbotapi.Message, productService *product.Service) {
	text := "Here all products: \n"
	for _, p := range productService.List() {
		text += p.Title
		text += "\n"

	}
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, text)
	_, _ = bot.Send(msg)

}

func helpCommand(bot *tgbotapi.BotAPI, inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "/help - help \n"+"/list - list entities")
	_, _ = bot.Send(msg)
}

func defaultBehavior(bot *tgbotapi.BotAPI, inputMessage *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "You wrote:"+inputMessage.Text)
	//msg.ReplyToMessageID = inputMessage.Message.MessageID
	_, _ = bot.Send(msg)
}
