package main

import (
	"binance-dca-bot-go/internal/config"
	"testing"

	"github.com/sudowanderer/notikit/notifier"
)

func TestBuildNotifierReturnsNilWithoutConfigs(t *testing.T) {
	notificationSender, err := buildNotifier(nil)
	if err != nil {
		t.Fatalf("buildNotifier returned error: %v", err)
	}
	if notificationSender != nil {
		t.Fatalf("expected nil notifier, got %T", notificationSender)
	}
}

func TestBuildNotifierBuildsBarkNotifier(t *testing.T) {
	notificationSender, err := buildNotifier([]config.NotificationConfig{
		{
			Type:      "bark",
			BarkKey:   "bark-key",
			BarkGroup: "binance-dca",
		},
	})
	if err != nil {
		t.Fatalf("buildNotifier returned error: %v", err)
	}

	barkNotifier, ok := notificationSender.(*notifier.BarkNotifier)
	if !ok {
		t.Fatalf("expected BarkNotifier, got %T", notificationSender)
	}
	if barkNotifier.Group != "binance-dca" {
		t.Fatalf("expected bark group %q, got %q", "binance-dca", barkNotifier.Group)
	}
}

func TestBuildNotifierBuildsMultiNotifier(t *testing.T) {
	notificationSender, err := buildNotifier([]config.NotificationConfig{
		{
			Type:      "bark",
			BarkKey:   "bark-key",
			BarkGroup: "binance-dca",
		},
		{
			Type:             "telegram",
			TelegramBotToken: "telegram-token",
			TelegramChatID:   "123456789",
		},
	})
	if err != nil {
		t.Fatalf("buildNotifier returned error: %v", err)
	}

	multiNotifier, ok := notificationSender.(*notifier.MultiNotifier)
	if !ok {
		t.Fatalf("expected MultiNotifier, got %T", notificationSender)
	}
	if len(multiNotifier.Notifiers) != 2 {
		t.Fatalf("expected 2 notifiers, got %d", len(multiNotifier.Notifiers))
	}
}
