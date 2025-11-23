use std::net::SocketAddr;
use std::sync::Arc;
use tokio::sync::mpsc;

use nostr_relay_builder::prelude::*;

#[derive(Debug, Clone)]
pub struct ReportNotification {
    pub event_id: String,
    pub message: String,
}

#[derive(Debug)]
pub struct ReportEvents {
    pub report_sender: Option<Arc<mpsc::UnboundedSender<ReportNotification>>>,
}

impl ReportEvents {
    pub fn new() -> Self {
        Self {
            report_sender: None,
        }
    }

    pub fn with_report_sender(
        mut self,
        sender: Arc<mpsc::UnboundedSender<ReportNotification>>,
    ) -> Self {
        self.report_sender = Some(sender);
        self
    }
}

impl WritePolicy for ReportEvents {
    fn admit_event<'a>(
        &'a self,
        event: &'a Event,
        addr: &'a SocketAddr,
    ) -> BoxedFuture<'a, PolicyResult> {
        Box::pin(async move {
            println!(
                "Checking event {:} with kind {:} (numeric: {}) from {}",
                event.id,
                event.kind,
                event.kind.as_u16(),
                addr
            );

            // Check if this is a reporting event (Kind 1984)
            // Using numeric comparison since Kind::Reporting might not be the right variant
            if event.kind.as_u16() == 1984 {
                log::info!("✅ Detected reporting event with kind 1984");
                log::info!(
                    "🚨 Received reporting event {} from {}",
                    event.id,
                    event.pubkey
                );

                // Extract report details (simplified parsing for now)
                let mut report_reason = "other".to_string();
                let mut reported_pubkey = None;
                let mut reported_event_id = None;
                let mut has_relevant_tags = false;

                // Try to find p or e tags by checking tag contents
                for tag in event.tags.iter() {
                    let tag_slice = tag.as_slice();
                    if tag_slice.len() >= 2 {
                        if tag_slice[0] == "p" {
                            reported_pubkey = Some(tag_slice[1].to_string());
                            if tag_slice.len() > 2 {
                                report_reason = tag_slice[2].to_string();
                            }
                            has_relevant_tags = true;
                        } else if tag_slice[0] == "e" {
                            reported_event_id = Some(tag_slice[1].to_string());
                            if tag_slice.len() > 2 {
                                report_reason = tag_slice[2].to_string();
                            }
                            has_relevant_tags = true;
                        }
                    }
                }

                if !has_relevant_tags {
                    log::info!("Ignoring report without p or e tags");
                    return PolicyResult::Accept;
                }

                // Format the report message
                let mut report_details = format!("🚨 **New Report (Kind 1984)**\n\n");
                report_details.push_str(&format!("**Report ID:** `{}`\n", event.id));
                report_details.push_str(&format!("**Reason:** {}\n", report_reason));
                report_details.push_str(&format!("**Reported by:** `{}`\n", event.pubkey));

                if !event.content.is_empty() {
                    report_details.push_str(&format!("**Additional info:** {}\n", event.content));
                }

                if let Some(pubkey) = &reported_pubkey {
                    report_details.push_str(&format!("**Reported pubkey:** `{}`\n", pubkey));
                }

                if let Some(event_id) = &reported_event_id {
                    report_details.push_str(&format!("**Reported event:** `{}`\n", event_id));
                }

                report_details.push_str("\n**Choose an action:**");

                log::info!("Report details: {}", report_details);

                // Send to Telegram bot if sender is available
                if let Some(sender) = &self.report_sender {
                    log::error!("1", );

                    let notification = ReportNotification {
                        event_id: event.id.to_string(),
                        message: report_details,
                    };
                    log::error!("2", );


                    if let Err(e) = sender.send(notification) {
                        log::error!("Failed to send report notification: {}", e);
                    } else {
                        log::info!("Report notification sent successfully");
                    }
                }
            }

            PolicyResult::Accept
        })
    }
}
