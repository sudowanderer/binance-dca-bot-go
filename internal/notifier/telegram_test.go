package notifier

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTelegramNotifier_Notify_Integration(t *testing.T) {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	// 初始化 TelegramNotifier
	notifier := &TelegramNotifier{
		BotToken: telegramBotToken,
		ChatID:   telegramChatID,
	}

	// 发送通知
	message := "Integration test: TelegramNotifier is working!"
	err := notifier.Notify(message)

	// 断言是否发送成功
	assert.NoError(t, err, "Notification should be sent successfully")
}
