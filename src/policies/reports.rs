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
        _: &'a SocketAddr,
    ) -> BoxedFuture<'a, PolicyResult> {
        Box::pin(async move {
            // TODO::: check the event more carefully to make sure it's not a spam
            // TODO::: pre-check the event with some moderation AI before sending to admins
            if event.kind.as_u16() == 1984 {
                log::info!(
                    "🚨 Received reporting event {} from {}",
                    event.id,
                    event.pubkey
                );

                let mut report_reason = "other".to_string();
                let mut reported_pubkey = None;
                let mut reported_event_id = None;
                let mut has_relevant_tags = false;

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

                let mut report_details = "🚨 New Report\n\n".to_string();
                report_details.push_str(&format!("Report Event ID: `{}`\n", event.id));
                report_details.push_str(&format!("Reason: {}\n", report_reason));
                report_details.push_str(&format!(
                    "Reported by: `https://npub.world/{}`\n",
                    event.pubkey
                ));

                if !event.content.is_empty() {
                    report_details.push_str(&format!("Additional info: {}\n", event.content));
                }

                if let Some(pubkey) = &reported_pubkey {
                    report_details.push_str(&format!(
                        "Reported pubkey: `https://npub.world/{}`\n",
                        pubkey
                    ));
                }

                if let Some(event_id) = &reported_event_id {
                    report_details.push_str(&format!("Reported event: `{}`\n", event_id));
                }

                report_details.push_str("\nChoose an action:");

                if let Some(sender) = &self.report_sender {
                    let notification = ReportNotification {
                        event_id: event.id.to_string(),
                        message: report_details,
                    };

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
