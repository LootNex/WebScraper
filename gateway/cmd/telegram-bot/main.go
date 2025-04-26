package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type loginResponse struct {
	Token string `json:"token"`
}

type registerResponse struct {
	UserID string `json:"user_id"`
}

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized as %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "login":
				handleLogin(update.Message, bot, update.Message.Chat.UserName)
			case "register":
				handleRegister(update.Message, bot, update.Message.Chat.UserName)
			case "logout":
				handleLogout(update.Message, bot, update.Message.Chat.UserName)
			case "check_item":
				handleCheckItem(update.Message, bot, update.Message.Chat.UserName)
			case "get_all_items":
				handleGetAllItems(update.Message, bot, update.Message.Chat.UserName)
			default:
				sendMessage(bot, update.Message.Chat.ID,
					"Unknown command. Try /login, /register, /logout, /check_item, /get_all_items")
			}
		}
	}
}

func handleLogin(message *tgbotapi.Message, bot *tgbotapi.BotAPI, telegramLogin string) {
	args := strings.Fields(message.CommandArguments())
	if len(args) != 2 {
		sendMessage(bot, message.Chat.ID, "Usage: /login <username> <password>")
		return
	}

	username, password := args[0], args[1]
	log.Printf("Received login: %s, password: %s", username, password)

	loginData := map[string]string{
		"login":          username,
		"password":       password,
		"telegram_login": telegramLogin,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to marshal request data")
		return
	}

	req, err := http.NewRequest("POST", os.Getenv("GATEWAY_URL")+"/login", bytes.NewBuffer(jsonData))
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
		return
	}
	defer resp.Body.Close()

	var loginResp loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		sendMessage(bot, message.Chat.ID, "Internal error")
		return
	}

	if loginResp.Token == "" {
		sendMessage(bot, message.Chat.ID, "Login failed")
		return
	}

	sendMessage(bot, message.Chat.ID, "Login successful")
}

func handleRegister(message *tgbotapi.Message, bot *tgbotapi.BotAPI, telegramLogin string) {
	args := strings.Fields(message.CommandArguments())
	if len(args) != 2 {
		sendMessage(bot, message.Chat.ID, "Usage: /register <username> <password>")
		return
	}

	username, password := args[0], args[1]
	log.Printf("Received login: %s, password: %s", username, password)

	registerData := map[string]string{
		"login":          username,
		"password":       password,
		"telegram_login": telegramLogin,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to marshal request data")
		return
	}

	req, err := http.NewRequest("POST", os.Getenv("GATEWAY_URL")+"/register", bytes.NewBuffer(jsonData))
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
		return
	}
	defer resp.Body.Close()

	var registerResp registerResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		sendMessage(bot, message.Chat.ID, "Internal error")
		return
	}

	if registerResp.UserID == "" {
		sendMessage(bot, message.Chat.ID, "Register failed")
		return
	}

	sendMessage(bot, message.Chat.ID, "Register successful, now you can log in")
}

func handleLogout(message *tgbotapi.Message, bot *tgbotapi.BotAPI, telegramLogin string) {
	args := strings.Fields(message.CommandArguments())
	if len(args) != 0 {
		sendMessage(bot, message.Chat.ID, "Usage: /logout")
		return
	}

	logoutData := map[string]string{
		"telegram_login": telegramLogin,
	}

	jsonData, err := json.Marshal(logoutData)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to marshal request data")
		return
	}

	req, err := http.NewRequest("POST", os.Getenv("GATEWAY_URL")+"/logout", bytes.NewBuffer(jsonData))
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
		return
	}
	defer resp.Body.Close()

	sendMessage(bot, message.Chat.ID, "Successful logout")
}

func handleCheckItem(message *tgbotapi.Message, bot *tgbotapi.BotAPI, telegramLogin string) {
	args := strings.Fields(message.CommandArguments())
	if len(args) != 1 {
		sendMessage(bot, message.Chat.ID, "Usage: /check_item <link>")
		return
	}

	link := args[0]
	log.Printf("Received link: %s", link)

	checkItemData := map[string]string{
		"telegram_login": telegramLogin,
		"link":           link,
	}

	jsonData, err := json.Marshal(checkItemData)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to marshal request data")
		return
	}

	req, _ := http.NewRequest("POST", os.Getenv("GATEWAY_URL")+"/check_item", bytes.NewBuffer(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
		return
	}
	defer resp.Body.Close()

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	if resp.StatusCode != http.StatusOK {
		sendMessage(bot, message.Chat.ID, "you need to login to use this method")
		return
	}

	sendMessage(bot, message.Chat.ID, fmt.Sprintf(
		"Item: %s\nStart Price: %.2f RUB\nCurrent Price: %.2f RUB\nDifference: %.2f",
		respBody["name"], respBody["start_price"], respBody["current_price"],
		respBody["difference_price"],
	))
}

func handleGetAllItems(message *tgbotapi.Message, bot *tgbotapi.BotAPI, telegramLogin string) {
	args := strings.Fields(message.CommandArguments())
	if len(args) != 0 {
		sendMessage(bot, message.Chat.ID, "Usage: /get_all_items")
		return
	}

	getAllItemsData := map[string]string{
		"telegram_login": telegramLogin,
	}

	jsonData, err := json.Marshal(getAllItemsData)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to marshal request data")
		return
	}

	req, _ := http.NewRequest("GET", os.Getenv("GATEWAY_URL")+"/get_all_items", bytes.NewBuffer(jsonData))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
		return
	}
	defer resp.Body.Close()

	var items []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&items)
	if resp.StatusCode != http.StatusOK {
		sendMessage(bot, message.Chat.ID, "you need to login to use this method")
		return
	}

	if len(items) == 0 {
		sendMessage(bot, message.Chat.ID, "No items found")
		return
	}

	var msg strings.Builder
	msg.WriteString("Your items:\n")
	for _, item := range items {
		msg.WriteString(fmt.Sprintf("- %s: %.2f RUB\n", item["name"], item["current_price"]))
	}
	sendMessage(bot, message.Chat.ID, msg.String())
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Failed to send message:", err)
	}
}
