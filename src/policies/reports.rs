use std::net::SocketAddr;

use nostr_relay_builder::prelude::*;

#[derive(Debug)]
pub struct ReportEvents;

impl WritePolicy for ReportEvents {
    fn admit_event<'a>(
        &'a self,
        event: &'a Event,
        addr: &'a SocketAddr,
    ) -> BoxedFuture<'a, PolicyResult> {
        Box::pin(async move {
            println!("Checking event {:} from {}", event.id, addr);

            if event.kind == Kind::Reporting {
                
            }
            PolicyResult::Accept
        })
    }
}
