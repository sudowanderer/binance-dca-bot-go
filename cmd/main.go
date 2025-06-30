package main

import (
	"binance-dca-bot-go/env"
	"binance-dca-bot-go/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	binanceconnector "github.com/binance/binance-connector-go"
	"github.com/sudowanderer/notikit/notifier"
	"strconv"
	"time"
)

func handleRequest(ctx context.Context, event json.RawMessage) error {
	// 1. 定义 provider 和 cfgEvent
	var provider config.ParameterProvider
	var cfgEvent config.ConfigEvent
	var err error

	if env.IsLambdaEnvironment() {
		// 云端：直接反序列化 EventBridge 传入的参数名字
		if err = json.Unmarshal(event, &cfgEvent); err != nil {
			return fmt.Errorf("failed to unmarshal event: %v", err)
		}
		// 用 SSM 拉取真正的参数值
		awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %v", err)
		}
		provider = &config.SSMParameterProvider{
			Client: ssm.NewFromConfig(awsCfg),
		}
	} else {
		// 本地：读本地 JSON 做 Mock
		var mock *config.LocalMockProvider
		mock, err = config.NewLocalMockProviderFromFile("local_config.json")
		if err != nil {
			return fmt.Errorf("failed to open local config: %v", err)
		}
		provider = mock
		// 把 mock.values 自动填到 cfgEvent 里
		if err = config.NewConfigEventFromLocalMock(mock, &cfgEvent); err != nil {
			return fmt.Errorf("failed to build config event from mock: %v", err)
		}
	}

	// 2. 统一加载 Config
	myCfg, err := config.LoadConfig(ctx, provider, cfgEvent)
	if err != nil {
		return fmt.Errorf("load config failed: %v", err)
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
		return fmt.Errorf("error placing order: %v", err)
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

func checkAndNotifyBalance(client *binanceconnector.Client, notifier notifier.Notifier, currency string, threshold *float64) error {
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
