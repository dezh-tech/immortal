pub mod types;

use std::fs;
use toml;
use crate::config::types::Config;

pub fn load(path: &str) -> Result<Config, Box<dyn std::error::Error>> {
    let s = fs::read_to_string(path)?;
    let cfg: Config = toml::from_str(&s)?;
    Ok(cfg)
}
