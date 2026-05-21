package events

import "time"

type CallStartedEvent struct {
	CallID    string    `json:"callId"`
	TenantID  string    `json:"tenantId"`
	FromUser  string    `json:"fromUser"`
	ToUser    string    `json:"toUser"`
	Direction string    `json:"direction"`
	StartedAt time.Time `json:"startedAt"`
}

type CallEndedEvent struct {
	CallID   string    `json:"callId"`
	TenantID string    `json:"tenantId"`
	EndedAt  time.Time `json:"endedAt"`
	Reason   string    `json:"reason"`
}
