# Binance DCA Bot

Binance DCA Bot is an open-source tool for automated Dollar-Cost Averaging (DCA) trading on Binance. This bot is designed to simplify cryptocurrency investments by automating periodic purchases of your target assets.

## Features

- **Automated DCA Orders**: Automatically place market buy orders for your target assets at predefined intervals.
- **Balance Monitoring**: Monitors your Binance account balance and sends warnings when the balance drops below a specified threshold.
- **Configurable Notifications**: Sends order failure and low-balance notifications through Bark, Telegram, or multiple channels.
- **AWS Parameter Store Support**: Fetches configuration parameters securely from AWS Parameter Store.

## Installation

### Prerequisites

- Go (version 1.20 or later)
- A Binance account with API key and secret
- A Bark key or Telegram bot for notifications
- AWS account with access to Parameter Store
- AWS SDK for Go

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/binance-dca-bot-go.git
   cd binance-dca-bot-go
   ```

2. Set up AWS Parameter Store with the required parameters:

   - `/myapp/bark/key`
   - `/myapp/binance/apiKey`
   - `/myapp/binance/apiSecret`

3. Build the binary:

   ```bash
   go build -o binance-dca-bot ./cmd/main.go
   ```

4. Run the bot:

   ```bash
   ./binance-dca-bot
   ```

For local testing, create personal config files from the checked-in examples:

```bash
cp local_config.example.json local_config.json
cp local_event.example.json local_event.json
```

Then fill `local_config.json` with your local secret values. The real local JSON files are ignored by Git.

## Usage

### Configuration

Configuration is managed through AWS Parameter Store. Ensure that the following parameters are set in your AWS account:

| Parameter Name                  | Description                                               |
| ------------------------------- | --------------------------------------------------------- |
| `/myapp/bark/key`               | Your Bark key for notifications.                          |
| `/myapp/binance/apiKey`         | Your Binance API key.                                     |
| `/myapp/binance/apiSecret`      | Your Binance API secret.                                  |

### Input Event Configuration

The bot requires an input event configuration in JSON format. Below is an example configuration:

```json
{
  "targetAssetParameter": "BTC",
  "amountParameter": "10",
  "orderCurrencyParameter": "FDUSD",
  "balanceThresholdParameter": "5000",
  "notifications": [
    {
      "type": "bark",
      "barkKeyParameter": "/myapp/bark/key",
      "barkGroupParameter": "binance-dca"
    }
  ],
  "binanceApiKeyParameter": "/myapp/binance/apiKey",
  "binanceApiSecretParameter": "/myapp/binance/apiSecret"
}
```

To send through Telegram instead:

```json
{
  "notifications": [
    {
      "type": "telegram",
      "telegramBotTokenParameter": "/myapp/telegram/token",
      "telegramChatIDParameter": "123456789"
    }
  ]
}
```

To send through multiple channels, include multiple entries in `notifications`. If `notifications` is not configured, the legacy `telegramBotTokenParameter` and `telegramChatIDParameter` fields are still supported.

### Notifications

The bot sends notifications to the configured channels when:

- An order fails.
- Your balance falls below the configured threshold.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push to your fork.
4. Submit a pull request with a detailed description of your changes.
