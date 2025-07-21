# Telegram Game Backend

## Quick Start
1. Copy example config:
   ```bash
   cp example.config.yaml config.yaml
   ```
2. Edit config.yaml with your credentials
3. Run the bot:
   ```bash
   go run main.go
   ```

## Configuration
Edit these fields in config.yaml:
- `telegram_token`: Get from @BotFather
- `game_short_name`: Your game's short name
- `port`: Server port (default: 8080)
- `game_url`: Where your game is hosted

A simple backend for a Telegram game built with Go.
