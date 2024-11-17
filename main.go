package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading dotenv: %s\n", err)
	}

	bot_token := os.Getenv("TELEGRAM_TOKEN")

	bot, err := telego.NewBot(bot_token)
	if err != nil {
		fmt.Println(err)
	}

	botUser, _ := bot.GetMe()
	fmt.Printf("Bot user: %+v\n", botUser)

	updates, _ := bot.UpdatesViaLongPolling(nil)

	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if update.Message != nil && (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") {
			chatID := update.Message.Chat.ID
			threadID := update.Message.MessageThreadID // For supergroups

			chatIdentifier := telego.ChatID{ID: chatID}
			admins, err := bot.GetChatAdministrators(&telego.GetChatAdministratorsParams{
				ChatID: chatIdentifier,
			})

			if err != nil {
				fmt.Println("Error getting admins:", err)
				return
			}

			var adminList []string
			for _, admin := range admins {
				if admin.MemberUser().Username != "caller_BDA_bot" && admin.MemberUser().Username != "" {
					adminList = append(adminList, "@"+admin.MemberUser().Username)
				}
			}

			messageText := phrases[rand.Intn(len(phrases))] + "\n" + fmt.Sprintf("%v", adminList)

			sendMessageParams := telego.SendMessageParams{
				ChatID:          chatIdentifier,
				Text:            messageText,
				MessageThreadID: threadID, // Include thread ID if it's part of a specific thread
			}

			_, _ = bot.SendMessage(&sendMessageParams)
		} else {
			chatID := update.Message.Chat.ID
			messageText := degenMsg

			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
		}
	}, th.CommandEqual("call"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if update.Message != nil && (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") {
			chatID := update.Message.Chat.ID

			chatIdentifier := telego.ChatID{ID: chatID}
			admins, err := bot.GetChatAdministrators(&telego.GetChatAdministratorsParams{
				ChatID: chatIdentifier,
			})

			if err != nil {
				fmt.Println("Error getting admins:", err)
				return
			}

			var adminList []string
			for _, admin := range admins {
				if admin.MemberUser().Username != "caller_BDA_bot" && admin.MemberUser().Username != "" {
					adminList = append(adminList, "@"+admin.MemberUser().Username)
				}
			}

			messageText := emergentPhrase[rand.Intn(len(emergentPhrase))] + "\n" + fmt.Sprintf("%v", adminList)
			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
		} else {
			chatID := update.Message.Chat.ID
			messageText := degenMsg

			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
		}
	}, th.CommandEqual("emergentCall"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if update.Message != nil && (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") {
			chatID := update.Message.Chat.ID

			messageText := leagueMsg
			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
		} else {
			chatID := update.Message.Chat.ID
			messageText := degenMsg

			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
		}
	}, th.CommandEqual("leaguecall"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := update.Message.Chat.ID

		messageText := helpMsg
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
	}, th.CommandEqual("help"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := update.Message.Chat.ID

		messageText := brokenMsg
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
	}, th.CommandEqual("broken"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		chatID := update.Message.Chat.ID
		messageText := startMessage
		_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
	}, th.CommandEqual("start"))

	bh.Start()
}
