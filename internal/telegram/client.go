package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func SendMessage(text string) error {
	token := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")
	threadID := os.Getenv("MESSAGE_THREAD_ID")
	if token == "" || chatID == "" {
		return fmt.Errorf("BOT_TOKEN or CHAT_ID not set")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}
	if threadID != "" {
		payload["message_thread_id"] = threadID
	}

	log.Println(token)
	log.Println(chatID)

	data, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Telegram error response: %s", string(body)) // 👈 Tambahkan baris ini
		return fmt.Errorf("Telegram error: %s", string(body))
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram error: %s", string(body))
	}

	return nil
}
