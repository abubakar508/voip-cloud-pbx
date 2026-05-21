package calls

import (
	"sync"
	"time"
)

type Direction string

const (
	DirectionInbound  Direction = "INBOUND"
	DirectionOutbound Direction = "OUTBOUND"
)

type CallSession struct {
	CallID     string
	TenantID   string
	FromUser   string
	ToUser     string
	Direction  Direction
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsFinished bool
}

type Manager struct {
	mu       sync.RWMutex
	sessions map[string]*CallSession
}

func NewManager() *Manager {
	return &Manager{
		sessions: make(map[string]*CallSession),
	}
}

func (m *Manager) Create(sess *CallSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[sess.CallID] = sess
}

func (m *Manager) Get(callID string) (*CallSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[callID]
	return s, ok
}

func (m *Manager) Finish(callID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[callID]; ok {
		s.IsFinished = true
		s.UpdatedAt = time.Now()
	}
}

func (m *Manager) Delete(callID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, callID)
}
