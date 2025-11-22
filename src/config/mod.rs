pub mod types;

use crate::config::types::Config;
use std::fs;
use toml;

pub fn load(path: &str) -> Result<Config, Box<dyn std::error::Error>> {
    let s = fs::read_to_string(path)?;
    let cfg: Config = toml::from_str(&s)?;
    Ok(cfg)
}
