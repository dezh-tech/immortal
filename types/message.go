package types

// Message repesents the messages between relays and clients.
type (
	Message     []string
	MessageType string
)

const (
	// Client to Relay.
	Request MessageType = "REQ"
	Close   MessageType = "CLOSE"

	// Relay to Client.
	Closed            MessageType = "CLOSED"
	Ok                MessageType = "OK"
	EndOFStoredEvents MessageType = "EOSE"
	Notice            MessageType = "NOTICE"

	// Both.
	Event MessageType = "EVENT"
)
