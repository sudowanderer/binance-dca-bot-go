package main

import (
	"binance-dca-bot-go/env"
	"binance-dca-bot-go/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	binanceconnector "github.com/binance/binance-connector-go"
	"github.com/sudowanderer/notikit/notifier"
)

func main() {
	if env.IsLambdaEnvironment() {
		// normal Lambda entrypoint
		lambda.Start(handleRequest)
		return
	}

	// --- local testing mode ---
	log.Println("🌱 Running in local mode, reading local_event.json …")

	data, err := os.ReadFile("local_event.json")
	if err != nil {
		log.Fatalf("failed to read event file: %v", err)
	}

	if err := handleRequest(context.Background(), data); err != nil {
		log.Fatalf("error in handleRequest: %v", err)
	}
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// parse event into ConfigEvent
	var cfgEvent config.ConfigEvent
	if err := json.Unmarshal(event, &cfgEvent); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// choose provider
	var provider config.ParameterProvider
	if env.IsLambdaEnvironment() {
		// AWS: real SSM
		awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
		if err != nil {
			return fmt.Errorf("loading AWS config: %w", err)
		}
		provider = &config.SSMParameterProvider{
			Client: ssm.NewFromConfig(awsCfg),
		}
	} else {
		// Local: mock from file
		mockProv, err := config.NewLocalMockProviderFromFile("local_config.json")
		if err != nil {
			return fmt.Errorf("failed to load local mock provider: %w", err)
		}
		provider = mockProv
	}

	// now call your existing loader (rename if needed)
	myCfg, err := config.LoadConfig(ctx, provider, cfgEvent)
	if err != nil {
		return fmt.Errorf("loading business config: %w", err)
	}

	if !env.IsLambdaEnvironment() {
		// ... the rest of your logic using myCfg ...
		log.Printf("Loaded config: %+v\n", myCfg)
	}

	// Optional timezone for Telegram messages
	loc, _ := time.LoadLocation("Asia/Shanghai")

	// Telegram Notifier
	tg := notifier.NewTelegramNotifierWithLocation(
		myCfg.TelegramBotToken,
		myCfg.TelegramChatID,
		loc,
	)

	// 初始化 Binance 客户端
	client := binanceconnector.NewClient(myCfg.BinanceAPIKey, myCfg.BinanceAPISecret)

	// 构造交易对符号
	symbol := myCfg.TargetAsset + myCfg.OrderCurrency

	// 下单
	newOrder, err := placeOrder(client, symbol, myCfg.Amount)
	if err != nil {
		errStrTemplate := "error placing order: %v"
		errStr := fmt.Sprintf(errStrTemplate, err)
		_ = tg.Notify(errStr)
		return fmt.Errorf(errStrTemplate, err)
	}
	fmt.Printf("Order placed: \n")
	fmt.Println(binanceconnector.PrettyPrint(newOrder))

	// 检查余额并发送通知
	err = checkAndNotifyBalance(client, tg, myCfg.OrderCurrency, myCfg.BalanceThreshold)
	if err != nil {
		return fmt.Errorf("error checking balance: %v", err)
	}

	return nil
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

func checkAndNotifyBalance(client *binanceconnector.Client, notifier notifier.Notifier, currency string, threshold *float64) error {
	balance, err := getBalance(client, currency)
	if err != nil {
		return fmt.Errorf("error fetching balance: %v", err)
	}

	balanceNum, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return fmt.Errorf("error parsing %s balance %q: %w", currency, balance, err)
	}
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
