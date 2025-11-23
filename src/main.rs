use nostr_lmdb::NostrLMDB;
use nostr_relay_builder::prelude::*;
use std::sync::Arc;
use tokio::sync::mpsc;

mod config;
mod policies;
mod manager;

use manager::TelegramBot;
use policies::reports::ReportNotification;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Initialize logging
    pretty_env_logger::init();
    log::info!("Starting Immortal...");

    let cfg = config::load("config.toml")?;
    let database = NostrLMDB::open(format!("{}/database", cfg.working_dir))?;
    
    // Create a second database instance for the bot (LMDB allows multiple readers)
    let bot_database = NostrLMDB::open(format!("{}/database", cfg.working_dir))?;
    let database_arc = Arc::new(bot_database);

    // Create channel for report notifications
    let (report_tx, mut report_rx) = mpsc::unbounded_channel::<ReportNotification>();
    let report_sender = Arc::new(report_tx);

    // Create report events policy with the sender
    let report_policy = policies::ReportEvents::new().with_report_sender(report_sender);
    
    let builder = RelayBuilder::default()
        .port(cfg.port)
        .database(database)
        .rate_limit(RateLimit {
            max_reqs: cfg.max_reqs,
            notes_per_minute: cfg.notes_per_minute,
        }).write_policy(report_policy);

    let relay = LocalRelay::new(builder);

    // Start relay
    relay.run().await?;
    log::info!("Relay listening on {}", relay.url().await);

    // Initialize Telegram bot if token is provided
    let bot_handle = if let Some(token) = &cfg.telegram_bot_token {
        log::info!("Starting Telegram bot...");
        let bot = Arc::new(TelegramBot::new(
            token.clone(), 
            cfg.telegram_group_chat_id, 
            database_arc.clone()
        ));
        
        let bot_for_commands = bot.clone();
        let bot_for_notifications = bot.clone();

        // Spawn command handler as a separate long-running task
        // Note: handle_commands runs the dispatcher indefinitely, so we don't wrap it in a Mutex
        let commands_handle = tokio::spawn(async move {
            log::info!("Telegram bot command handler starting...");
            if let Err(e) = bot_for_commands.handle_commands().await {
                log::error!("Telegram bot error: {}", e);
            }
        });

        // Spawn notification handler as a separate concurrent task
        let notifications_handle = tokio::spawn(async move {
            log::info!("Notification handler started, waiting for reports...");
            while let Some(notification) = report_rx.recv().await {
                log::info!("📨 Received report notification for event: {}", notification.event_id);
                
                // Send report notification directly (no lock needed since send_report_notification takes &self)
                match bot_for_notifications.send_report_notification(&notification.event_id, &notification.message).await {
                    Ok(()) => {
                        log::info!("✅ Successfully sent Telegram notification");
                    }
                    Err(e) => {
                        log::error!("❌ Failed to send Telegram notification: {}", e);
                    }
                }
            }
            log::warn!("Notification handler loop ended!");
        });

        Some((commands_handle, notifications_handle))
    } else {
        log::info!("No Telegram bot token provided, bot will not start");
        None
    };

    // Wait for Ctrl+C signal
    tokio::signal::ctrl_c().await?;

    log::info!("Shutting down...");

    // Cancel bot tasks if they exist
    if let Some((commands_handle, notifications_handle)) = bot_handle {
        commands_handle.abort();
        notifications_handle.abort();
    }

    Ok(())
}
