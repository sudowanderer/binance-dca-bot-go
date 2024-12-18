package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Notifier interface {
	Notify(message string) error
}

type TelegramNotifier struct {
	BotToken string
	ChatID   string
}

func (t *TelegramNotifier) Notify(message string) error {
	// 获取当前时间（上海时区）
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %v", err)
	}
	currentTime := time.Now().In(location).Format("2006-01-02 15:04:05")

	// 拼接消息和时间戳
	fullMessage := fmt.Sprintf("%s\n%s", message, currentTime)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)
	payload := map[string]string{
		"chat_id": t.ChatID,
		"text":    fullMessage,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
