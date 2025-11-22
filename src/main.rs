use nostr_lmdb::NostrLMDB;
use nostr_relay_builder::prelude::*;

mod config;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let cfg = config::load("config.toml")?; 
    let database = NostrLMDB::open(format!("{}/database", cfg.working_dir))?;
    let builder = RelayBuilder::default()
        .port(cfg.port)
        .database(database)
        .rate_limit(RateLimit {
            max_reqs: 128,
            notes_per_minute: 30,
        });

    let relay = LocalRelay::new(builder);

    relay.run().await?;

    println!("Relay listening on {}", relay.url().await);

    tokio::signal::ctrl_c().await?;

    Ok(())
}
