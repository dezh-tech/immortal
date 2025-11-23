use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub working_dir: String,
    pub port: u16,
    pub max_reqs: usize,
    pub notes_per_minute: u32,
    pub telegram_bot_token: Option<String>,
    pub telegram_group_chat_id: Option<i64>,
}
