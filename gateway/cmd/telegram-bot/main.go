package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/joho/godotenv"
)

const gatewayURL = "http://api-gateway:8080" // URL вашего API Gateway

type loginResponse struct {
    Token   string `json:"token"`
    Message string `json:"message"`
}

var userTokens = make(map[int64]string) // Храним токены пользователей

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

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
                handleLogin(update.Message, bot)
            default:
                sendMessage(bot, update.Message.Chat.ID, "Unknown command. Try /login, /additem, /checkitem, /getallitems")
            }
        }
    }
}

func handleLogin(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
    args := strings.Fields(message.CommandArguments())
    if len(args) != 2 {
        sendMessage(bot, message.Chat.ID, "Usage: /login <username> <password>")
        return
    }

    username, password := args[0], args[1]
    log.Printf("Received login: %s, password: %s", username, password) 

    loginData := map[string]string{
        "login":    username,
        "password": password,
    }

    jsonData, err := json.Marshal(loginData)
    if err != nil {
        sendMessage(bot, message.Chat.ID, "Failed to marshal request data")
        return
    }

    req, err := http.NewRequest("POST", gatewayURL+"/login", bytes.NewBuffer(jsonData))
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
        sendMessage(bot, message.Chat.ID, "Login failed: "+loginResp.Message)
        return
    }

    userTokens[message.From.ID] = loginResp.Token
    sendMessage(bot, message.Chat.ID, "Login successful")
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
    msg := tgbotapi.NewMessage(chatID, text)
    _, err := bot.Send(msg)
    if err != nil {
        log.Println("Failed to send message:", err)
    }
}



// func handleLogin(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
//  	args := strings.Fields(message.CommandArguments())
//  	if len(args) != 2 {
//   		sendMessage(bot, message.Chat.ID, "Usage: /login <username> <password>")
//   	return
// 	}
// 	log.Printf("Received login: %s, password: %s", args[0], args[1]) // Отладка
//  	resp, err := http.PostForm(gatewayURL+"/login", map[string][]string{
// 		"login":    {args[0]},
// 		"password": {args[1]},
// 	})
// 	if err != nil {
// 		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
// 		return
// 	}
//  	defer resp.Body.Close()

// 	var loginResp loginResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
// 		sendMessage(bot, message.Chat.ID, "Internal error")
// 		return
// 	}
 

//  	if !loginResp.Success {
//   		sendMessage(bot, message.Chat.ID, loginResp.Message)
//   		return
// 	}

//  	userTokens[message.From.ID] = loginResp.Token
//  	sendMessage(bot, message.Chat.ID, "Login successful")
// }

// func handleAddItem(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
// 	token, ok := userTokens[message.From.ID]
// 	if !ok {
// 		sendMessage(bot, message.Chat.ID, "Please login first with /login")
// 		return
// 	}

// 	args := strings.TrimSpace(message.CommandArguments())
// 	if args == "" {
// 		sendMessage(bot, message.Chat.ID, "Usage: /additem <link>")
// 		return
// 	}

// 	reqBody, _ := json.Marshal(map[string]string{"link": args})
// 	req, _ := http.NewRequest("POST", gatewayURL+"/additem", bytes.NewBuffer(reqBody))
// 	req.Header.Set("Authorization", token)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
//   		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
// 		return
// 	}
//  	defer resp.Body.Close()

//  	var respBody map[string]interface{}
//  	json.NewDecoder(resp.Body).Decode(&respBody)
// 	if resp.StatusCode != http.StatusOK {
//   		sendMessage(bot, message.Chat.ID, fmt.Sprintf("Error: %v", respBody["message"]))
// 		return
// 	}

// 	sendMessage(bot, message.Chat.ID, fmt.Sprintf("%v", respBody["message"]))
// }

// func handleCheckItem(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
//  	token, ok := userTokens[message.From.ID]
//  	if !ok {
//   		sendMessage(bot, message.Chat.ID, "Please login first with /login")
//  		 return
// 	}

//  	args := strings.TrimSpace(message.CommandArguments())
// 	if args == "" {
//   		sendMessage(bot, message.Chat.ID, "Usage: /checkitem <link>")
//  		return
// 	}

//  	req, _ := http.NewRequest("GET", gatewayURL+"/checkitem?link="+args, nil)
//  	req.Header.Set("Authorization", token)

//  	client := &http.Client{}
//  	resp, err := client.Do(req)
//  	if err != nil {
//  		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
//   	return
// 	}
//  	defer resp.Body.Close()

// 	var respBody map[string]interface{}
// 	json.NewDecoder(resp.Body).Decode(&respBody)
// 	if resp.StatusCode != http.StatusOK {
// 		sendMessage(bot, message.Chat.ID, fmt.Sprintf("Error: %v", respBody["message"]))
// 		return
// 	}

// 	sendMessage(bot, message.Chat.ID, fmt.Sprintf(
// 		"Item: %s\nPrice: %.2f RUB\nLast checked: %s",
// 		respBody["link"], respBody["current_price"], respBody["last_checked"],
// 	))
// }

// func handleGetAllItems(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
//  	token, ok := userTokens[message.From.ID]
//  	if !ok {
// 		sendMessage(bot, message.Chat.ID, "Please login first with /login")
//   		return
// 	}

// 	req, _ := http.NewRequest("GET", gatewayURL+"/getallitems", nil)
// 	req.Header.Set("Authorization", token)

//  	client := &http.Client{}
//  	resp, err := client.Do(req)
//  	if err != nil {
//   		sendMessage(bot, message.Chat.ID, "Failed to connect to server")
// 		return
// 	}
// 	defer resp.Body.Close()

//  	var respBody map[string]interface{}
// 	json.NewDecoder(resp.Body).Decode(&respBody)
//  	if resp.StatusCode != http.StatusOK {
// 		sendMessage(bot, message.Chat.ID, fmt.Sprintf("Error: %v", respBody["message"]))
// 		return
// 	}

//  	items := respBody["items"].([]interface{})
//  	if len(items) == 0 {
// 		sendMessage(bot, message.Chat.ID, "No items found")
//   		return
// 	}

//  	var msg strings.Builder
//  	msg.WriteString("Your items:\n")
//  	for _, item := range items {
//   		i := item.(map[string]interface{})
//   		msg.WriteString(fmt.Sprintf("- %s: %.2f RUB\n", i["link"], i["current_price"]))
// 	}
// 	sendMessage(bot, message.Chat.ID, msg.String())
// }

// func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
//  	msg := tgbotapi.NewMessage(chatID, text)
//  	bot.Send(msg)
// }