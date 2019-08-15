package hooks

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

var telegramToken, _ = os.LookupEnv("TELEGRAM_TOKEN")
var telegramChannel, _ = os.LookupEnv("TELEGRAM_CHANNEL")
var telegramapi = "https://api.telegram.org/bot" + telegramToken + "/sendMessage?chat_id=" + telegramChannel + "&text="

// SendTelegramMessage sends telegram message with message
func SendTelegramMessage(message string) {
	fmt.Println("Sending message to telegram")
	url := telegramapi + url.QueryEscape(message)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send message to telegram, %v", err)
		os.Exit(1)
	}
}
