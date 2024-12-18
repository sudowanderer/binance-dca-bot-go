package main

import (
	"binance-dca-bot-go/internal/config"
	mynotifier "binance-dca-bot-go/internal/notifier"
	"context"
	"fmt"
	binanceconnector "github.com/binance/binance-connector-go"
	"os"
	"strconv"
)

func main() {
	envConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config: %+v\n", envConfig)

	var notifier mynotifier.Notifier
	if envConfig.TelegramChatID != "" && envConfig.TelegramBotToken != "" {
		notifier = &mynotifier.TelegramNotifier{
			ChatID:   envConfig.TelegramChatID,
			BotToken: envConfig.TelegramBotToken,
		}
	}

	// Initialise the client
	client := binanceconnector.NewClient(envConfig.BinanceAPIKey, envConfig.BinanceAPISecret)

	// Get symbol
	symbol := envConfig.TargetAsset + envConfig.OrderCurrency

	// Create new order
	newOrder, err := placeOrder(client, symbol, envConfig.Amount)
	handleError(err, "Error placing order")
	fmt.Printf("Order placed: \n")
	fmt.Println(binanceconnector.PrettyPrint(newOrder))

	// check balances
	err = checkAndNotifyBalance(client, notifier, envConfig.OrderCurrency, envConfig.BalanceThreshold)
	handleError(err, "Error checking balance")

}

func getBalance(client *binanceconnector.Client, asset string) (string, error) {
	accountInfo, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return "", fmt.Errorf("error fetching balance: %v", err)
	}

	for _, balance := range accountInfo.Balances {
		if balance.Asset == asset {
			return balance.Free, nil
		}
	}

	return "0", nil
}

func placeOrder(client *binanceconnector.Client, symbol string, amount float64) (any, error) {
	return client.NewCreateOrderService().Symbol(symbol).
		Side("BUY").Type("MARKET").QuoteOrderQty(amount).
		Do(context.Background())
}

func checkAndNotifyBalance(client *binanceconnector.Client, notifier mynotifier.Notifier, currency string, threshold *float64) error {
	balance, err := getBalance(client, currency)
	if err != nil {
		return fmt.Errorf("error fetching balance: %v", err)
	}

	balanceNum, _ := strconv.ParseFloat(balance, 64)
	if threshold != nil && balanceNum < *threshold {
		message := fmt.Sprintf("Warning: Your %s balance is below the threshold of %.2f. Current balance: %.2f", currency, *threshold, balanceNum)
		if notifier != nil {
			return notifier.Notify(message)
		}
	}
	return nil
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
		os.Exit(1)
	}
}
