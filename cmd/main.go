package main

import (
	"binance-dca-bot-go/internal/config"
	mynotifier "binance-dca-bot-go/internal/notifier"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	binanceconnector "github.com/binance/binance-connector-go"
	"strconv"
)

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// 加载配置
	envConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	fmt.Printf("Config: %+v\n", envConfig)

	// 初始化通知器
	var notifier mynotifier.Notifier
	if envConfig.TelegramChatID != "" && envConfig.TelegramBotToken != "" {
		notifier = &mynotifier.TelegramNotifier{
			ChatID:   envConfig.TelegramChatID,
			BotToken: envConfig.TelegramBotToken,
		}
	}

	// 初始化 Binance 客户端
	client := binanceconnector.NewClient(envConfig.BinanceAPIKey, envConfig.BinanceAPISecret)

	// 构造交易对符号
	symbol := envConfig.TargetAsset + envConfig.OrderCurrency

	// 下单
	newOrder, err := placeOrder(client, symbol, envConfig.Amount)
	if err != nil {
		return fmt.Errorf("error placing order: %v", err)
	}
	fmt.Printf("Order placed: \n")
	fmt.Println(binanceconnector.PrettyPrint(newOrder))

	// 检查余额并发送通知
	err = checkAndNotifyBalance(client, notifier, envConfig.OrderCurrency, envConfig.BalanceThreshold)
	if err != nil {
		return fmt.Errorf("error checking balance: %v", err)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
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
			if err := notifier.Notify(message); err != nil {
				return fmt.Errorf("error sending notification: %v", err)
			}
		}
	}
	return nil
}
