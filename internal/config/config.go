package config

import (
	"context"
	"fmt"
	"strconv"
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

// ParameterProvider 是参数加载的抽象接口
type ParameterProvider interface {
	GetParameter(ctx context.Context, paramName string) (string, error)
}

// LoadConfig 从 provider 加载所有参数
func LoadConfig(ctx context.Context, provider ParameterProvider, event ConfigEvent) (*Config, error) {
	targetAsset, err := provider.GetParameter(ctx, event.TargetAssetParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get TARGET_ASSET: %v", err)
	}

	amountStr, err := provider.GetParameter(ctx, event.AmountParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get AMOUNT: %v", err)
	}

	orderCurrency, err := provider.GetParameter(ctx, event.OrderCurrencyParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get ORDER_CURRENCY: %v", err)
	}

	balanceThresholdStr, err := provider.GetParameter(ctx, event.BalanceThresholdParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get BALANCE_THRESHOLD: %v", err)
	}

	telegramBotToken, err := provider.GetParameter(ctx, event.TelegramBotTokenParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get TELEGRAM_BOT_TOKEN: %v", err)
	}

	telegramChatID, err := provider.GetParameter(ctx, event.TelegramChatIDParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get TELEGRAM_CHAT_ID: %v", err)
	}

	binanceAPIKey, err := provider.GetParameter(ctx, event.BinanceApiKeyParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get BINANCE_API_KEY: %v", err)
	}

	binanceAPISecret, err := provider.GetParameter(ctx, event.BinanceApiSecretParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get BINANCE_API_SECRET: %v", err)
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("AMOUNT must be a valid number: %v", err)
	}

	var balanceThreshold *float64
	if balanceThresholdStr != "" {
		threshold, err := strconv.ParseFloat(balanceThresholdStr, 64)
		if err != nil {
			return nil, fmt.Errorf("BALANCE_THRESHOLD must be a valid number: %v", err)
		}
		balanceThreshold = &threshold
	}

	return &Config{
		TargetAsset:      targetAsset,
		Amount:           amount,
		OrderCurrency:    orderCurrency,
		BalanceThreshold: balanceThreshold,
		TelegramBotToken: telegramBotToken,
		TelegramChatID:   telegramChatID,
		BinanceAPIKey:    binanceAPIKey,
		BinanceAPISecret: binanceAPISecret,
	}, nil
}
