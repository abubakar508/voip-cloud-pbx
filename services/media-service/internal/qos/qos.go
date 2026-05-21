package qos

import (
	"sync"

	"github.com/pion/rtp"
)

type StreamKey struct {
	SSRC uint32
	Addr string // remote address string (ip:port)
}

type StreamStats struct {
	Packets     uint64 `json:"packets"`
	Lost        uint64 `json:"lost"`
	LastSeq     uint16 `json:"lastSeq"`
	Initialized bool   `json:"-"`
}

type Manager struct {
	mu      sync.RWMutex
	streams map[StreamKey]*StreamStats
}

func NewManager() *Manager {
	return &Manager{
		streams: make(map[StreamKey]*StreamStats),
	}
}

// UpdateFromPacket updates QoS stats for a given stream key with a new RTP packet.
func (m *Manager) UpdateFromPacket(key StreamKey, pkt *rtp.Packet) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats, ok := m.streams[key]
	if !ok {
		stats = &StreamStats{}
		m.streams[key] = stats
	}

	if !stats.Initialized {
		stats.LastSeq = pkt.SequenceNumber
		stats.Packets = 1
		stats.Initialized = true
		return
	}

	// Calculate loss based on sequence number gaps (basic)
	expectedNext := stats.LastSeq + 1
	if pkt.SequenceNumber != expectedNext {
		// Handle wrap-around (16-bit sequence)
		if pkt.SequenceNumber > expectedNext {
			stats.Lost += uint64(pkt.SequenceNumber - expectedNext)
		} else if expectedNext > pkt.SequenceNumber {
			// Wrap around scenario, approximate loss
			stats.Lost += uint64((65535 - expectedNext) + pkt.SequenceNumber + 1)
		}
	}

	stats.Packets++
	stats.LastSeq = pkt.SequenceNumber
}

// Snapshot returns a copy of all QoS stats for exposure via HTTP.
func (m *Manager) Snapshot() map[StreamKey]StreamStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[StreamKey]StreamStats, len(m.streams))
	for k, v := range m.streams {
		out[k] = *v
	}
	return out
}
