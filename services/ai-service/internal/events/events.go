package events

import "time"

type CallEndedEvent struct {
	CallID   string    `json:"callId"`
	TenantID string    `json:"tenantId"`
	EndedAt  time.Time `json:"endedAt"`
	Reason   string    `json:"reason"`
}
