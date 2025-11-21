use nostr_lmdb::NostrLMDB;
use nostr_relay_builder::prelude::*;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Open a database (all databases that implements `NostrDatabase` trait can be used).
    let database = NostrLMDB::open("nostr-relay")?;
    // Configure the relay.
    let builder = RelayBuilder::default()
        .port(7777)
        .database(database)
        .rate_limit(RateLimit {
            max_reqs: 128,
            notes_per_minute: 30,
        });

    // Construct the relay instance.
    let relay = LocalRelay::new(builder);

    // Start the relay.
    relay.run().await?;

    println!("Relay listening on {}", relay.url().await);

    // Keep the process running.
    tokio::signal::ctrl_c().await?;

    Ok(())
}
