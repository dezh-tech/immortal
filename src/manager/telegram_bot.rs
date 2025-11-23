use teloxide::prelude::*;
use teloxide::utils::command::BotCommands;
use std::error::Error;

pub struct TelegramBot {
    bot: Bot,
}

impl TelegramBot {
    pub fn new(token: String) -> Self {
        let bot = Bot::new(token);
        Self { bot }
    }

    pub async fn start(&self) -> Result<(), Box<dyn Error + Send + Sync>> {
        log::info!("Starting Telegram bot...");
        
        teloxide::repl(self.bot.clone(), |bot: Bot, msg: Message| async move {
            bot.send_message(msg.chat.id, "Hello from Immortal bot!").await?;
            Ok(())
        })
        .await;
        
        Ok(())
    }

    pub async fn send_notification(&self, chat_id: i64, message: &str) -> Result<(), Box<dyn Error + Send + Sync>> {
        self.bot.send_message(ChatId(chat_id), message).await?;
        Ok(())
    }

    pub async fn handle_commands(&self) -> Result<(), Box<dyn Error + Send + Sync>> {
        let handler = dptree::entry()
            .branch(
                Update::filter_message()
                    .branch(
                        dptree::entry()
                            .filter_command::<Command>()
                            .endpoint(Self::command_handler),
                    )
                    .branch(dptree::endpoint(Self::message_handler)),
            );

        Dispatcher::builder(self.bot.clone(), handler)
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
