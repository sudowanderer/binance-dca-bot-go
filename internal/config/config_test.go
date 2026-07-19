package config

import (
	"context"
	"testing"
)

type stubParameterProvider map[string]string

func (p stubParameterProvider) GetParameter(ctx context.Context, name string) (string, error) {
	value, ok := p[name]
	if !ok {
		return "", &missingParameterError{name: name}
	}
	return value, nil
}

type missingParameterError struct {
	name string
}

func (e *missingParameterError) Error() string {
	return "missing parameter: " + e.name
}

func TestLoadConfigLoadsNotificationList(t *testing.T) {
	cfg, err := LoadConfig(context.Background(), stubParameterProvider{
		"/myapp/bark/key":          "bark-key",
		"/myapp/telegram/token":    "telegram-token",
		"/myapp/binance/apiKey":    "binance-key",
		"/myapp/binance/apiSecret": "binance-secret",
	}, ConfigEvent{
		TargetAssetParam:      "BTC",
		AmountParam:           "10",
		OrderCurrencyParam:    "FDUSD",
		BalanceThresholdParam: "5000",
		BinanceApiKeyParam:    "/myapp/binance/apiKey",
		BinanceApiSecretParam: "/myapp/binance/apiSecret",
		Notifications: []NotificationEvent{
			{
				Type:           " Bark ",
				BarkKeyParam:   "/myapp/bark/key",
				BarkGroupParam: "binance-dca",
			},
			{
				Type:                  "telegram",
				TelegramBotTokenParam: "/myapp/telegram/token",
				TelegramChatIDParam:   "123456789",
			},
		},
	})
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if len(cfg.Notifications) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(cfg.Notifications))
	}
	if cfg.Notifications[0].Type != "bark" || cfg.Notifications[0].BarkKey != "bark-key" || cfg.Notifications[0].BarkGroup != "binance-dca" {
		t.Fatalf("unexpected bark notification config: %+v", cfg.Notifications[0])
	}
	if cfg.Notifications[1].Type != "telegram" || cfg.Notifications[1].TelegramBotToken != "telegram-token" || cfg.Notifications[1].TelegramChatID != "123456789" {
		t.Fatalf("unexpected telegram notification config: %+v", cfg.Notifications[1])
	}
}

func TestLoadConfigFallsBackToLegacyTelegramFields(t *testing.T) {
	cfg, err := LoadConfig(context.Background(), stubParameterProvider{
		"/myapp/telegram/token":    "telegram-token",
		"/myapp/binance/apiKey":    "binance-key",
		"/myapp/binance/apiSecret": "binance-secret",
	}, ConfigEvent{
		TargetAssetParam:      "BTC",
		AmountParam:           "10",
		OrderCurrencyParam:    "FDUSD",
		BinanceApiKeyParam:    "/myapp/binance/apiKey",
		BinanceApiSecretParam: "/myapp/binance/apiSecret",
		TelegramBotTokenParam: "/myapp/telegram/token",
		TelegramChatIDParam:   "123456789",
	})
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if len(cfg.Notifications) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(cfg.Notifications))
	}
	if cfg.Notifications[0].Type != "telegram" || cfg.Notifications[0].TelegramBotToken != "telegram-token" || cfg.Notifications[0].TelegramChatID != "123456789" {
		t.Fatalf("unexpected legacy telegram notification config: %+v", cfg.Notifications[0])
	}
}

func TestLoadConfigReturnsErrorForUnsupportedNotificationType(t *testing.T) {
	_, err := LoadConfig(context.Background(), stubParameterProvider{
		"/myapp/binance/apiKey":    "binance-key",
		"/myapp/binance/apiSecret": "binance-secret",
	}, ConfigEvent{
		TargetAssetParam:      "BTC",
		AmountParam:           "10",
		OrderCurrencyParam:    "FDUSD",
		BinanceApiKeyParam:    "/myapp/binance/apiKey",
		BinanceApiSecretParam: "/myapp/binance/apiSecret",
		Notifications: []NotificationEvent{
			{Type: "email"},
		},
	})
	if err == nil {
		t.Fatal("expected unsupported notification type error")
	}
}
