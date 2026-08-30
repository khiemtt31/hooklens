package inbox

import "time"

// InboxEvent represents one HTTP request captured by an inbox.
type InboxEvent struct {
	ID         int64               `json:"id"`
	InboxID    int64               `json:"inboxId"`
	Method     string              `json:"method"`
	Headers    map[string][]string `json:"headers"`
	Body       string              `json:"body"`
	ReceivedAt time.Time           `json:"receivedAt"`
}
