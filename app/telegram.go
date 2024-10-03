package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

// SendToTelegram sends a message to the specified chat ID using a Telegram bot.
func SendToTelegram(
	chatID string, 
	token string, 
	request *http.Request, 
	blacklistedActions string,
	blacklistedEvents string,
	) string {
	messageText, err := formatGitHubWebhook(request, blacklistedActions, blacklistedEvents)
	if err != nil {
		return "nothing to send"
	}

	tgStatus := "failed"
	if SendMessage(chatID, messageText, token) {
		tgStatus = "succeed"
	}
	return tgStatus
}

// SendMessage sends a message to Telegram API.
func SendMessage(chatID, text string, token string) bool {
	apiURL := "https://api.telegram.org/bot" + token + "/sendMessage"
	data := map[string]interface{}{
		"chat_id":                    chatID,
		"text":                       text,
		"parse_mode":                 "HTML",
		"disable_web_page_preview":   true,
		"disable_notification":       true,
	}

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error sending message to Telegram:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var tgResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&tgResponse); err == nil {
			if msgLink := GetTelegramMessageLink(tgResponse); msgLink != "" {
				log.Println("Telegram message link:", msgLink)
				return true
			}
		}
	}
	return false
}

// GetTelegramMessageLink retrieves the message link from the Telegram response.
func GetTelegramMessageLink(jsonData map[string]interface{}) string {
	result, ok := jsonData["result"].(map[string]interface{})
	if !ok {
		log.Println("No result for link")
		return ""
	}

	messageID := int(result["message_id"].(float64))
	chatID := int(result["chat"].(map[string]interface{})["id"].(float64))
	msgLink := fmt.Sprintf("https://t.me/%d/%d", chatID, messageID)
	return msgLink
}
