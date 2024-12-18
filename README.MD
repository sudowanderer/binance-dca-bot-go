# Binance DCA Bot

Binance DCA Bot is an open-source tool for automated Dollar-Cost Averaging (DCA) trading on Binance. This bot is designed to simplify cryptocurrency investments by automating periodic purchases of your target assets.

## Features

- **Automated DCA Orders**: Automatically place market buy orders for your target assets at predefined intervals.
- **Balance Monitoring**: Monitors your Binance account balance and sends warnings when the balance drops below a specified threshold.
- **Telegram Notifications**: Sends real-time notifications to your Telegram account for order confirmations and low-balance warnings.

## Installation

### Prerequisites

- Go (version 1.20 or later)
- A Binance account with API key and secret
- A Telegram bot for notifications

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/binance-dca-bot-go.git
   cd binance-dca-bot-go
   ```

2. Set up your environment variables:

   ```bash
   TARGET_ASSET=BTC
   ORDER_CURRENCY=USDT
   AMOUNT=50
   BALANCE_THRESHOLD=100
   TELEGRAM_BOT_TOKEN=your_telegram_bot_token
   TELEGRAM_CHAT_ID=your_telegram_chat_id
   BINANCE_API_KEY=your_binance_api_key
   BINANCE_API_SECRET=your_binance_api_secret
   ```

3. Build the binary:

   ```bash
   go build -o binance_dca_bot ./cmd/main.go
   ```

4. Run the bot:

   ```bash
   ./binance-dca-bot
   ```

## Usage

### Configuration

All configuration is managed through environment variables. Use the `.env` file or export variables directly in your shell.

| Variable             | Description                                               |
| -------------------- | --------------------------------------------------------- |
| `TARGET_ASSET`       | The asset you want to buy (e.g., `BTC`).                  |
| `ORDER_CURRENCY`     | The currency used to buy the target asset (e.g., `USDT`). |
| `AMOUNT`             | The amount to spend per DCA order.                        |
| `BALANCE_THRESHOLD`  | The balance threshold to trigger a low-balance alert.     |
| `TELEGRAM_BOT_TOKEN` | Your Telegram bot token for notifications.                |
| `TELEGRAM_CHAT_ID`   | Your Telegram chat ID for notifications.                  |
| `BINANCE_API_KEY`    | Your Binance API key.                                     |
| `BINANCE_API_SECRET` | Your Binance API secret.                                  |

### Notifications

The bot sends notifications to your Telegram account when:

- An order is placed.
- Your balance falls below the configured threshold.

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push to your fork.
4. Submit a pull request with a detailed description of your changes.