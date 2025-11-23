use nostr_lmdb::NostrLMDB;
use nostr_relay_builder::prelude::*;

mod config;
mod policies;
mod manager;

use manager::TelegramBot;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize logging
    pretty_env_logger::init();
    log::info!("Starting Immortal...");

    let cfg = config::load("config.toml")?;
    let database = NostrLMDB::open(format!("{}/database", cfg.working_dir))?;
    let builder = RelayBuilder::default()
        .port(cfg.port)
        .database(database)
        .rate_limit(RateLimit {
            max_reqs: cfg.max_reqs,
            notes_per_minute: cfg.notes_per_minute,
        }).write_policy(policies::ReportEvents);

    let relay = LocalRelay::new(builder);

    // Start relay
    relay.run().await?;
    log::info!("Relay listening on {}", relay.url().await);

    // Initialize Telegram bot if token is provided
    let bot_handle = if let Some(token) = &cfg.telegram_bot_token {
        log::info!("Starting Telegram bot...");
        let bot = TelegramBot::new(token.clone());
        Some(tokio::spawn(async move {
            if let Err(e) = bot.handle_commands().await {
                log::error!("Telegram bot error: {}", e);
            }
        }))
    } else {
        log::info!("No Telegram bot token provided, bot will not start");
        None
    };

    // Wait for Ctrl+C signal
    tokio::signal::ctrl_c().await?;

    log::info!("Shutting down...");

    // Cancel bot task if it exists
    if let Some(handle) = bot_handle {
        handle.abort();
    }

    Ok(())
}
