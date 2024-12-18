package config

import (
	"fmt"
	"os"
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

func LoadConfig() (*Config, error) {
	targetAsset := os.Getenv("TARGET_ASSET")
	amountStr := os.Getenv("AMOUNT")
	orderCurrency := os.Getenv("ORDER_CURRENCY")
	balanceThresholdStr := os.Getenv("BALANCE_THRESHOLD")
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	binanceAPIKey := os.Getenv("BINANCE_API_KEY")
	binanceAPISecret := os.Getenv("BINANCE_API_SECRET")

	if targetAsset == "" || amountStr == "" || orderCurrency == "" {
		return nil, fmt.Errorf("environment variables TARGET_ASSET, AMOUNT, and ORDER_CURRENCY must be set")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("AMOUNT environment variable must be a valid number: %v", err)
	}

	var balanceThreshold *float64
	if balanceThresholdStr != "" {
		threshold, err := strconv.ParseFloat(balanceThresholdStr, 64)
		if err != nil {
			return nil, fmt.Errorf("BALANCE_THRESHOLD environment variable must be a valid number: %v", err)
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
