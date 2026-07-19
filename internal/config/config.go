package config

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	TargetAsset      string
	Amount           float64
	OrderCurrency    string
	BalanceThreshold *float64
	Notifications    []NotificationConfig
	BinanceAPIKey    string
	BinanceAPISecret string
}

type NotificationConfig struct {
	Type             string
	TelegramBotToken string
	TelegramChatID   string
	BarkKey          string
	BarkGroup        string
}

// ConfigEvent 表示 EventBridge 传入的参数 key
type ConfigEvent struct {
	TargetAssetParam      string              `json:"targetAssetParameter"`
	AmountParam           string              `json:"amountParameter"`
	OrderCurrencyParam    string              `json:"orderCurrencyParameter"`
	BalanceThresholdParam string              `json:"balanceThresholdParameter"`
	Notifications         []NotificationEvent `json:"notifications"`

	// Deprecated: use Notifications instead.
	TelegramBotTokenParam string `json:"telegramBotTokenParameter"`
	// Deprecated: use Notifications instead.
	TelegramChatIDParam   string `json:"telegramChatIDParameter"`
	BinanceApiKeyParam    string `json:"binanceApiKeyParameter"`
	BinanceApiSecretParam string `json:"binanceApiSecretParameter"`
}

type NotificationEvent struct {
	Type                  string `json:"type"`
	TelegramBotTokenParam string `json:"telegramBotTokenParameter"`
	TelegramChatIDParam   string `json:"telegramChatIDParameter"`
	BarkKeyParam          string `json:"barkKeyParameter"`
	BarkGroupParam        string `json:"barkGroupParameter"`
}

func getMaybeParam(ctx context.Context, provider ParameterProvider, raw string) (string, error) {
	if strings.HasPrefix(raw, "/") {
		return provider.GetParameter(ctx, raw)
	}
	return raw, nil
}

func LoadConfig(ctx context.Context, provider ParameterProvider, event ConfigEvent) (*Config, error) {
	// for non-secrets, these will just return the literal "BTC", "100.0", etc.
	targetAsset, err := getMaybeParam(ctx, provider, event.TargetAssetParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get TARGET_ASSET: %w", err)
	}

	amountStr, err := getMaybeParam(ctx, provider, event.AmountParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get AMOUNT: %w", err)
	}

	orderCurrency, err := getMaybeParam(ctx, provider, event.OrderCurrencyParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get ORDER_CURRENCY: %w", err)
	}

	// BalanceThreshold might be empty string, you can similarly treat that
	thRaw, err := getMaybeParam(ctx, provider, event.BalanceThresholdParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get BALANCE_THRESHOLD: %w", err)
	}

	binanceAPIKey, err := provider.GetParameter(ctx, event.BinanceApiKeyParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get BINANCE_API_KEY: %w", err)
	}
	binanceAPISecret, err := provider.GetParameter(ctx, event.BinanceApiSecretParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get BINANCE_API_SECRET: %w", err)
	}

	// parse the floats as before
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("AMOUNT must be a valid number: %w", err)
	}
	var balanceThreshold *float64
	if thRaw != "" {
		t, err := strconv.ParseFloat(thRaw, 64)
		if err != nil {
			return nil, fmt.Errorf("BALANCE_THRESHOLD must be a valid number: %w", err)
		}
		balanceThreshold = &t
	}

	notifications, err := loadNotifications(ctx, provider, event)
	if err != nil {
		return nil, err
	}

	return &Config{
		TargetAsset:      targetAsset,
		Amount:           amount,
		OrderCurrency:    orderCurrency,
		BalanceThreshold: balanceThreshold,
		Notifications:    notifications,
		BinanceAPIKey:    binanceAPIKey,
		BinanceAPISecret: binanceAPISecret,
	}, nil
}

func loadNotifications(ctx context.Context, provider ParameterProvider, event ConfigEvent) ([]NotificationConfig, error) {
	notificationEvents := event.Notifications
	if len(notificationEvents) == 0 && event.TelegramBotTokenParam != "" && event.TelegramChatIDParam != "" {
		notificationEvents = []NotificationEvent{
			{
				Type:                  "telegram",
				TelegramBotTokenParam: event.TelegramBotTokenParam,
				TelegramChatIDParam:   event.TelegramChatIDParam,
			},
		}
	}

	notifications := make([]NotificationConfig, 0, len(notificationEvents))
	for _, notificationEvent := range notificationEvents {
		notification, err := loadNotification(ctx, provider, notificationEvent)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func loadNotification(ctx context.Context, provider ParameterProvider, event NotificationEvent) (NotificationConfig, error) {
	notificationType := strings.ToLower(strings.TrimSpace(event.Type))
	switch notificationType {
	case "telegram":
		telegramBotToken, err := provider.GetParameter(ctx, event.TelegramBotTokenParam)
		if err != nil {
			return NotificationConfig{}, fmt.Errorf("failed to get TELEGRAM_BOT_TOKEN: %w", err)
		}
		telegramChatID, err := getMaybeParam(ctx, provider, event.TelegramChatIDParam)
		if err != nil {
			return NotificationConfig{}, fmt.Errorf("failed to get TELEGRAM_CHAT_ID: %w", err)
		}
		return NotificationConfig{
			Type:             notificationType,
			TelegramBotToken: telegramBotToken,
			TelegramChatID:   telegramChatID,
		}, nil
	case "bark":
		barkKey, err := provider.GetParameter(ctx, event.BarkKeyParam)
		if err != nil {
			return NotificationConfig{}, fmt.Errorf("failed to get BARK_KEY: %w", err)
		}
		barkGroup, err := getMaybeParam(ctx, provider, event.BarkGroupParam)
		if err != nil {
			return NotificationConfig{}, fmt.Errorf("failed to get BARK_GROUP: %w", err)
		}
		return NotificationConfig{
			Type:      notificationType,
			BarkKey:   barkKey,
			BarkGroup: barkGroup,
		}, nil
	default:
		return NotificationConfig{}, fmt.Errorf("unsupported notification type %q", event.Type)
	}
}
