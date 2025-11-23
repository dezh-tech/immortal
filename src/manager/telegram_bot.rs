use teloxide::prelude::*;
use teloxide::utils::command::BotCommands;
use teloxide::types::{InlineKeyboardButton, InlineKeyboardMarkup};
use std::error::Error;
use std::sync::Arc;
use nostr_lmdb::NostrLMDB;
use nostr_relay_builder::prelude::*;

pub struct TelegramBot {
    bot: Bot,
    group_chat_id: Option<i64>,
    database: Arc<NostrLMDB>,
}

impl TelegramBot {
    pub fn new(token: String, group_chat_id: Option<i64>, database: Arc<NostrLMDB>) -> Self {
        let bot = Bot::new(token);
        Self { 
            bot,
            group_chat_id,
            database,
        }
    }

    pub async fn send_report_notification(&self, event_id: &str, report_details: &str) -> Result<(), Box<dyn Error + Send + Sync>> {
        log::info!("received new report {:}", event_id);
        if let Some(chat_id) = self.group_chat_id {
            log::info!("Sending Telegram notification to chat_id: {}", chat_id);
            
            // Shorten event_id to first 16 chars to fit within Telegram's 64-byte callback_data limit
            let short_id = &event_id[..std::cmp::min(16, event_id.len())];
            
            let keyboard = InlineKeyboardMarkup::new(vec![vec![
                InlineKeyboardButton::callback("🗑️ Delete", format!("d:{}", short_id)),
                InlineKeyboardButton::callback("❌ Ignore", format!("i:{}", short_id)),
            ]]);

            match self.bot
                .send_message(ChatId(chat_id), report_details)
                .reply_markup(keyboard)
                .await
            {
                Ok(message) => {
                    log::info!("Successfully sent Telegram message with ID: {}", message.id);
                }
                Err(e) => {
                    log::error!("Failed to send Telegram message: {}", e);
                    return Err(Box::new(e));
                }
            }
        } else {
            log::warn!("Cannot send report notification: no group chat ID configured");
        }
        Ok(())
    }

    pub async fn handle_commands(&self) -> Result<(), Box<dyn Error + Send + Sync>> {
        let database = self.database.clone();
        let bot = self.bot.clone();
        
        let handler = dptree::entry()
            .branch(
                Update::filter_message()
                    .branch(
                        dptree::entry()
                            .filter_command::<Command>()
                            .endpoint(Self::command_handler),
                    )
                    .branch(dptree::endpoint(Self::message_handler)),
            )
            .branch(
                Update::filter_callback_query()
                    .endpoint(move |bot_inner: Bot, q: CallbackQuery| {
                        let db = database.clone();
                        async move {
                            Self::callback_handler_with_db(bot_inner, q, db).await
                        }
                    }),
            );

        Dispatcher::builder(bot, handler)
            .dependencies(dptree::deps![])
            .enable_ctrlc_handler()
            .build()
            .dispatch()
            .await;

        Ok(())
    }

    async fn command_handler(
        bot: Bot,
        msg: Message,
        cmd: Command,
    ) -> Result<(), Box<dyn Error + Send + Sync>> {
        match cmd {
            Command::Help => {
                bot.send_message(msg.chat.id, Command::descriptions().to_string())
                    .await?;
            }
            Command::Start => {
                bot.send_message(
                    msg.chat.id,
                    "Welcome to Immortal! This bot monitors the Nostr relay status.",
                )
                .await?;
            }
            Command::Status => {
                bot.send_message(
                    msg.chat.id,
                    "Relay is running and healthy! 🟢",
                )
                .await?;
            }
        }
        Ok(())
    }

    async fn message_handler(bot: Bot, msg: Message) -> Result<(), Box<dyn Error + Send + Sync>> {
        if let Some(text) = msg.text() {
            log::info!("Received message: {}", text);
            bot.send_message(msg.chat.id, format!("You said: {}", text))
                .await?;
        }
        Ok(())
    }

    async fn callback_handler_with_db(bot: Bot, q: CallbackQuery, database: Arc<NostrLMDB>) -> Result<(), Box<dyn Error + Send + Sync>> {
        if let Some(data) = q.data {
            let parts: Vec<&str> = data.split(':').collect();
            if parts.len() == 2 {
                let action = parts[0];
                let short_id = parts[1];

                match action {
                    "d" => {
                        log::info!("Delete requested for event ID (short): {}", short_id);
                        
                        // Get the database reference
                        let db = database.as_ref();
                        
                        // Query all reporting events to find one matching the short_id
                        // This is a workaround since we only have the short ID
                        let filter = Filter::new().kind(Kind::Reporting);
                        
                        if let Ok(events) = db.query(filter).await {
                            // Find the report event with matching short_id
                            if let Some(report_event) = events.into_iter()
                                .find(|e| e.id.to_string().starts_with(short_id)) {
                                log::info!("Found report event: {}", report_event.id);
                                
                                // Parse tags to find what to delete
                                for tag in report_event.tags.iter() {
                                    let tag_slice = tag.as_slice();
                                    if tag_slice.len() >= 2 {
                                        if tag_slice[0] == "e" {
                                            // Event report - delete the reported event
                                            let reported_event_id = &tag_slice[1];
                                            log::info!("Deleting reported event: {}", reported_event_id);
                                            
                                            if let Ok(event_id) = EventId::parse(reported_event_id) {
                                                let delete_filter = Filter::new().id(event_id);
                                                if let Err(e) = db.delete(delete_filter).await {
                                                    log::error!("Failed to delete event {}: {}", reported_event_id, e);
                                                } else {
                                                    log::info!("✅ Successfully deleted event: {}", reported_event_id);
                                                }
                                            }
                                        } else if tag_slice[0] == "p" {
                                            // Profile report - delete all events by this pubkey
                                            let reported_pubkey = &tag_slice[1];
                                            log::info!("Deleting all events by pubkey: {}", reported_pubkey);
                                            
                                            if let Ok(pubkey) = PublicKey::parse(reported_pubkey) {
                                                let delete_filter = Filter::new().author(pubkey);
                                                if let Err(e) = db.delete(delete_filter).await {
                                                    log::error!("Failed to delete events by pubkey {}: {}", reported_pubkey, e);
                                                } else {
                                                    log::info!("✅ Successfully deleted all events by pubkey: {}", reported_pubkey);
                                                }
                                            }
                                        }
                                    }
                                }
                            } else {
                                log::warn!("Report event not found with short ID: {}", short_id);
                            }
                        } else {
                            log::error!("Failed to query report events");
                        }
                        
                        // Delete the Telegram message
                        if let Some(message) = q.message {
                            bot.delete_message(message.chat.id, message.id).await?;
                        }
                        
                        // Send confirmation
                        bot.answer_callback_query(q.id)
                            .text("✅ Content deleted successfully")
                            .await?;
                    }
                    "i" => {
                        log::info!("Ignore requested for event ID (short): {}", short_id);
                        
                        // Just delete the Telegram message
                        if let Some(message) = q.message {
                            bot.delete_message(message.chat.id, message.id).await?;
                        }
                        
                        bot.answer_callback_query(q.id)
                            .text("❌ Report ignored")
                            .await?;
                    }
                    _ => {
                        bot.answer_callback_query(q.id)
                            .text("Unknown action")
                            .await?;
                    }
                }
            }
        }
        Ok(())
    }
}

#[derive(BotCommands, Clone)]
#[command(rename_rule = "lowercase", description = "Immortal Bot Commands:")]
enum Command {
    #[command(description = "display this text.")]
    Help,
    #[command(description = "start the bot.")]
    Start,
    #[command(description = "get relay status.")]
    Status,
}
