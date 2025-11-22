use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub working_dir: String,
    pub port: u16,
}

