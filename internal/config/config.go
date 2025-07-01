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
	TelegramBotToken string
	TelegramChatID   string
	BinanceAPIKey    string
	BinanceAPISecret string
}

// ConfigEvent 表示 EventBridge 传入的参数 key
type ConfigEvent struct {
	TargetAssetParam      string `json:"targetAssetParameter"`
	AmountParam           string `json:"amountParameter"`
	OrderCurrencyParam    string `json:"orderCurrencyParameter"`
	BalanceThresholdParam string `json:"balanceThresholdParameter"`
	TelegramBotTokenParam string `json:"telegramBotTokenParameter"`
	TelegramChatIDParam   string `json:"telegramChatIDParameter"`
	BinanceApiKeyParam    string `json:"binanceApiKeyParameter"`
	BinanceApiSecretParam string `json:"binanceApiSecretParameter"`
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

	telegramChatID, err := getMaybeParam(ctx, provider, event.TelegramChatIDParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get TELEGRAM_CHAT_ID: %w", err)
	}

	// BalanceThreshold might be empty string, you can similarly treat that
	thRaw, err := getMaybeParam(ctx, provider, event.BalanceThresholdParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get BALANCE_THRESHOLD: %w", err)
	}

	// _Only_ for the truly secret ones do we force a GetParameter
	telegramBotToken, err := provider.GetParameter(ctx, event.TelegramBotTokenParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get TELEGRAM_BOT_TOKEN: %w", err)
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

	return &Config{
		TargetAsset:      targetAsset,
		Amount:           amount,
		OrderCurrency:    orderCurrency,
		BalanceThreshold: balanceThreshold,
		TelegramBotToken: telegramBotToken,
		TelegramChatID:   telegramChatID, // assuming you GetParameter’d that above
		BinanceAPIKey:    binanceAPIKey,
		BinanceAPISecret: binanceAPISecret,
	}, nil
}
