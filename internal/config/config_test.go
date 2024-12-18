package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("TARGET_ASSET", "BTC")
	os.Setenv("AMOUNT", "0.01")
	os.Setenv("ORDER_CURRENCY", "USDT")
	os.Setenv("BALANCE_THRESHOLD", "100.5")
	os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")
	os.Setenv("TELEGRAM_CHAT_ID", "12345")
	os.Setenv("BINANCE_API_KEY", "api_key")
	os.Setenv("BINANCE_API_SECRET", "api_secret")

	defer func() {
		// Clean up environment variables
		os.Unsetenv("TARGET_ASSET")
		os.Unsetenv("AMOUNT")
		os.Unsetenv("ORDER_CURRENCY")
		os.Unsetenv("BALANCE_THRESHOLD")
		os.Unsetenv("TELEGRAM_BOT_TOKEN")
		os.Unsetenv("TELEGRAM_CHAT_ID")
		os.Unsetenv("BINANCE_API_KEY")
		os.Unsetenv("BINANCE_API_SECRET")
	}()

	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "BTC", config.TargetAsset)
	assert.Equal(t, 0.01, config.Amount)
	assert.Equal(t, "USDT", config.OrderCurrency)
	assert.NotNil(t, config.BalanceThreshold)
	assert.Equal(t, 100.5, *config.BalanceThreshold)
	assert.Equal(t, "test_token", config.TelegramBotToken)
	assert.Equal(t, "12345", config.TelegramChatID)
	assert.Equal(t, "api_key", config.BinanceAPIKey)
	assert.Equal(t, "api_secret", config.BinanceAPISecret)
}

func TestLoadConfig_MissingRequiredVariables(t *testing.T) {
	// Unset required environment variables
	os.Unsetenv("TARGET_ASSET")
	os.Unsetenv("AMOUNT")
	os.Unsetenv("ORDER_CURRENCY")

	config, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestLoadConfig_InvalidAmount(t *testing.T) {
	// Set up invalid amount
	os.Setenv("TARGET_ASSET", "BTC")
	os.Setenv("AMOUNT", "invalid")
	os.Setenv("ORDER_CURRENCY", "USDT")

	defer func() {
		os.Unsetenv("TARGET_ASSET")
		os.Unsetenv("AMOUNT")
		os.Unsetenv("ORDER_CURRENCY")
	}()

	config, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
}
