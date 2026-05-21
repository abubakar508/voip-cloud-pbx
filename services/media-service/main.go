package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"github.com/nats-io/nats.go"
	"github.com/pion/rtp"
	"go.uber.org/zap"

	"github.com/abubakar508/voip-cloud-pbx/services/media-service/internal/calls"
	httpHandlers "github.com/abubakar508/voip-cloud-pbx/services/media-service/internal/http"
	"github.com/abubakar508/voip-cloud-pbx/services/media-service/internal/natsclient"
	"github.com/abubakar508/voip-cloud-pbx/services/media-service/internal/qos"
)

func main() {
	bootstrap := shared.Init("media-service")
	shared.PrintBanner("media-service")

	logger := bootstrap.Logger

	httpPort := envOrDefault("MEDIA_SERVICE_PORT", "8082")
	rtpPort := envOrInt("MEDIA_RTP_BASE_PORT", 40000)

	httpAddr := ":" + httpPort
	rtpAddr := fmt.Sprintf("0.0.0.0:%d", rtpPort)

	// Managers
	qm := qos.NewManager()
	callMgr := calls.NewManager()

	// HTTP server and routes
	server := httpserver.New(httpserver.Options{
		Addr:   httpAddr,
		Logger: logger,
	})
	engine := server.Engine()

	qosHandler := httpHandlers.NewQoSHandler(qm)
	qosHandler.RegisterRoutes(engine)

	callsHandler := httpHandlers.NewCallsHandler(callMgr)
	callsHandler.RegisterRoutes(engine)

	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("media http server error", zap.Error(err))
		}
	}()

	// NATS client (unchanged, but reuse callMgr from above)
	nc, err := natsclient.New(bootstrap.Config)
	if err != nil {
		logger.Fatal("failed to connect to nats", zap.Error(err))
	}
	defer nc.Close()
	logger.Info("media-service connected to NATS", zap.String("url", nc.URL()))

	// Subscribe to call events
	if _, err := nc.Subscribe("calls.started", func(msg *nats.Msg) {
		var evt calls.CallStartedEvent
		if err := nc.Unmarshal(msg, &evt); err != nil {
			logger.Error("failed to unmarshal calls.started", zap.Error(err))
			return
		}
		sess := &calls.CallSession{
			CallID:    evt.CallID,
			TenantID:  evt.TenantID,
			FromUser:  evt.FromUser,
			ToUser:    evt.ToUser,
			Direction: calls.CallDirection(evt.Direction),
			StartedAt: evt.StartedAt,
		}
		callMgr.Upsert(sess)
		logger.Info("registered call session from NATS (started)",
			zap.String("callId", evt.CallID),
			zap.String("tenantId", evt.TenantID),
			zap.String("from", evt.FromUser),
			zap.String("to", evt.ToUser),
		)
	}); err != nil {
		logger.Fatal("failed to subscribe to calls.started", zap.Error(err))
	}

	if _, err := nc.Subscribe("calls.ended", func(msg *nats.Msg) {
		var evt calls.CallEndedEvent
		if err := nc.Unmarshal(msg, &evt); err != nil {
			logger.Error("failed to unmarshal calls.ended", zap.Error(err))
			return
		}
		callMgr.Finish(evt.CallID, evt.EndedAt)
		logger.Info("updated call session from NATS (ended)",
			zap.String("callId", evt.CallID),
			zap.String("tenantId", evt.TenantID),
			zap.String("reason", evt.Reason),
		)
	}); err != nil {
		logger.Fatal("failed to subscribe to calls.ended", zap.Error(err))
	}

	// Register routes for the calls list handler
	callsListHandler := httpHandlers.NewCallsListHandler(callMgr)
	callsListHandler.RegisterRoutes(engine)

	// RTP UDP listener
	udpAddr, err := net.ResolveUDPAddr("udp", rtpAddr)
	if err != nil {
		logger.Fatal("failed to resolve udp addr", zap.Error(err))
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Fatal("failed to listen on udp", zap.Error(err))
	}
	logger.Info("Media RTP UDP listener started",
		zap.String("addr", rtpAddr),
	)

	buf := make([]byte, 1500)
	packet := &rtp.Packet{}

	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Error("error reading udp", zap.Error(err))
			continue
		}
		raw := buf[:n]

		if err := packet.Unmarshal(raw); err != nil {
			logger.Warn("failed to parse rtp packet", zap.Error(err))
			continue
		}

		// Update QoS
		key := qos.StreamKey{
			SSRC: packet.SSRC,
			Addr: remote.String(),
		}
		qm.UpdateFromPacket(key, packet)

		// Forwarding: find call session where this addr is A or B
		if sess, leg, ok := callMgr.FindSessionByAddr(remote); ok {
			var target *net.UDPAddr
			if leg == "A" && sess.B.Addr != nil {
				target = sess.B.Addr
			} else if leg == "B" && sess.A.Addr != nil {
				target = sess.A.Addr
			}

			if target != nil {
				if _, err := conn.WriteToUDP(raw, target); err != nil {
					logger.Error("failed to forward RTP packet",
						zap.String("callId", sess.CallID),
						zap.String("leg", leg),
						zap.String("from", remote.String()),
						zap.String("to", target.String()),
						zap.Error(err),
					)
				} else {
					logger.Debug("forwarded RTP packet",
						zap.String("callId", sess.CallID),
						zap.String("leg", leg),
						zap.String("from", remote.String()),
						zap.String("to", target.String()),
					)
				}
			}
		}

		logger.Info("RTP packet received",
			zap.String("from", remote.String()),
			zap.Uint32("ssrc", packet.SSRC),
			zap.Uint16("seq", packet.SequenceNumber),
			zap.Uint32("timestamp", packet.Timestamp),
			zap.Int("payload_type", int(packet.PayloadType)),
			zap.Int("payload_len", len(packet.Payload)),
			zap.Bool("marker", packet.Marker),
		)
	}
}

func envOrDefault(name, def string) string {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	return v
}

func envOrInt(name string, def int) int {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
