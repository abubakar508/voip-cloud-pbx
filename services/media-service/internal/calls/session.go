package calls

import (
	"net"
	"sync"
	"time"
)

type CallDirection string

const (
	DirectionInbound  CallDirection = "INBOUND"
	DirectionOutbound CallDirection = "OUTBOUND"
)

type RTPFlow struct {
	Addr *net.UDPAddr
	SSRC uint32
}

type CallSession struct {
	CallID    string
	TenantID  string
	FromUser  string
	ToUser    string
	Direction CallDirection
	StartedAt time.Time
	EndedAt   *time.Time

	// RTP endpoints for each leg
	A RTPFlow
	B RTPFlow
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

func (m *Manager) Upsert(sess *CallSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[sess.CallID] = sess
}

func (m *Manager) Finish(callID string, endedAt time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[callID]; ok {
		s.EndedAt = &endedAt
	}
}

func (m *Manager) Delete(callID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, callID)
}

func (m *Manager) Get(callID string) (*CallSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[callID]
	return s, ok
}

// SetEndpoints sets the RTP endpoints for call legs A and B.
func (m *Manager) SetEndpoints(callID string, aAddr *net.UDPAddr, bAddr *net.UDPAddr) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[callID]; ok {
		s.A.Addr = aAddr
		s.B.Addr = bAddr
	}
}

// FindSessionByAddr returns the call session and a flag indicating if the given remote address matches A or B leg.
func (m *Manager) FindSessionByAddr(addr *net.UDPAddr) (*CallSession, string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, s := range m.sessions {
		if s.A.Addr != nil && s.A.Addr.IP.Equal(addr.IP) && s.A.Addr.Port == addr.Port {
			return s, "A", true
		}
		if s.B.Addr != nil && s.B.Addr.IP.Equal(addr.IP) && s.B.Addr.Port == addr.Port {
			return s, "B", true
		}
	}
	return nil, "", false
}

func (m *Manager) Snapshot() []*CallSession {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]*CallSession, 0, len(m.sessions))
	for _, s := range m.sessions {
		out = append(out, s)
	}
	return out
}
