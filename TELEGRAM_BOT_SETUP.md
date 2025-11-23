# Telegram Bot Setup Guide

This guide will help you set up the Telegram bot integration for Immortal.

## Prerequisites

1. A Telegram account
2. Access to @BotFather on Telegram

## Step 1: Create a Telegram Bot

1. Open Telegram and search for `@BotFather`
2. Start a conversation with BotFather by clicking "Start" or sending `/start`
3. Create a new bot by sending `/newbot`
4. Follow the prompts:
   - Choose a name for your bot (e.g., "Immortal Relay Monitor")
   - Choose a username for your bot (e.g., "immortal_relay_bot")
5. BotFather will provide you with a token that looks like: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyZ`
6. **Save this token securely** - you'll need it for configuration

## Step 2: Configure Immortal

You have two options to configure the bot token:

### Option A: Using config.toml (Development)

1. Open your `config.toml` file
2. Add the bot token:
   ```toml
   telegram_bot_token = "YOUR_BOT_TOKEN_HERE"
   ```
   Replace `YOUR_BOT_TOKEN_HERE` with the actual token from BotFather

### Option B: Using Environment Variable (Recommended for Production)

1. Set the environment variable:
   ```bash
   export TELEGRAM_BOT_TOKEN="YOUR_BOT_TOKEN_HERE"
   ```

Note: Environment variable takes precedence over config.toml if both are set.

## Step 3: Run Immortal with Bot

1. Start Immortal:
   ```bash
   cargo run --release
   ```
   
2. You should see logs indicating the bot is starting:
   ```
   INFO  immortal > Starting Immortal...
   INFO  immortal > Relay listening on ws://127.0.0.1:7777
   INFO  immortal > Starting Telegram bot...
   ```

## Step 4: Test the Bot

1. Find your bot on Telegram using the username you created
2. Start a conversation with your bot by clicking "Start" or sending `/start`
3. Try these commands:
   - `/start` - Get a welcome message
   - `/help` - See available commands
   - `/status` - Check relay status

## Available Bot Commands

- `/start` - Welcome message and bot introduction
- `/help` - Display list of available commands
- `/status` - Get current relay status

## Troubleshooting

### Bot doesn't respond

1. **Check the token**: Ensure you copied the entire token correctly
2. **Check logs**: Look for error messages in the application logs
3. **Environment**: Make sure `RUST_LOG=info` is set to see debug information

### Bot shows offline

1. **Application running**: Ensure Immortal is running
2. **Network connectivity**: Check if the server has internet access
3. **Token validity**: Verify the token is still valid with BotFather

### Example Error Messages

- `No Telegram bot token provided` - Token not configured in config.toml
- `Telegram bot error: ...` - Check the specific error message for details

## Security Notes

- Keep your bot token private and secure
- Don't commit the token to version control
- Consider using environment variables for production deployments
- The bot token should be treated like a password

## Adding More Features

The bot system is extensible. You can add more commands by:

1. Adding new variants to the `Command` enum in `src/manager/telegram_bot.rs`
2. Implementing handlers in the `command_handler` function
3. Adding descriptions to help users understand the commands

## Production Deployment

For production use, consider:

- Using environment variables instead of config files for sensitive data
- Setting up proper logging and monitoring
- Implementing rate limiting for bot commands
- Adding admin-only commands for sensitive operations
