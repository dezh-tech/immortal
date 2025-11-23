pub mod types;

use crate::config::types::Config;
use std::fs;

pub fn load(path: &str) -> Result<Config, Box<dyn std::error::Error>> {
    let s = fs::read_to_string(path)?;
    let mut cfg: Config = toml::from_str(&s)?;

    // Override with environment variable if provided
    if let Ok(token) = std::env::var("TELEGRAM_BOT_TOKEN") {
        cfg.telegram_bot_token = Some(token);
    }

    Ok(cfg)
}
