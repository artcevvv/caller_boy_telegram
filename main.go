package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading dotenv: %s\n", err)
	}

	bot_token := os.Getenv("TELEGRAM_TOKEN")

	// Create a new bot instance
	bot, err := telego.NewBot(bot_token)
	if err != nil {
		fmt.Println(err)
	}

	botUser, _ := bot.GetMe()
	fmt.Printf("Bot user: %+v\n", botUser)

	// Get updates via long polling
	updates, _ := bot.UpdatesViaLongPolling(nil)

	// Create bot handler
	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	// Handle the "call" command
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Only process group or supergroup chats
		if update.Message != nil && (update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup") {
			chatID := update.Message.Chat.ID

			// Get chat administrators for the current chat
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

			// Send the message with the list of admins in the current chat
			messageText := "Calling everyone... \n" + fmt.Sprintf("%v", adminList)
			_, _ = bot.SendMessage(tu.Message(tu.ID(chatID), messageText))
		}
	}, th.CommandEqual("call"))

	// Start the bot handler
	bh.Start()
}
