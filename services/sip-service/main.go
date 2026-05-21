package main

import (
	"context"
	"net"
	"os"
	"time"

	shared "github.com/abubakar508/voip-cloud-pbx/packages/shared-go"
	"github.com/abubakar508/voip-cloud-pbx/packages/shared-go/httpserver"
	"github.com/ghettovoice/gosip/sip"
	"github.com/ghettovoice/gosip/sip/parser"
	"go.uber.org/zap"

	"github.com/abubakar508/voip-cloud-pbx/services/sip-service/internal/auth"
	"github.com/abubakar508/voip-cloud-pbx/services/sip-service/internal/calls"
	"github.com/abubakar508/voip-cloud-pbx/services/sip-service/internal/models"
	"github.com/abubakar508/voip-cloud-pbx/services/sip-service/internal/natsclient"
	"github.com/abubakar508/voip-cloud-pbx/services/sip-service/internal/redisclient"
	"github.com/abubakar508/voip-cloud-pbx/services/sip-service/internal/registration"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	bootstrap := shared.Init("sip-service")
	shared.PrintBanner("sip-service")

	logger := bootstrap.Logger

	// Postgres: auto-migrate SipAccount
	gdb, err := gorm.Open(postgres.Open(bootstrap.Config.PostgresDSN), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}
	if err := gdb.AutoMigrate(&models.SipAccount{}); err != nil {
		logger.Fatal("failed to auto-migrate sip models", zap.Error(err))
	}

	authSvc := auth.NewService(gdb)

	// Redis client and registration store
	redisCli := redisclient.New(bootstrap.Config)
	if err := redisCli.Ping(context.Background()); err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	regStore := registration.NewRedisStore(redisCli.SetWithTTL)

	// NATS client
	nc, err := natsclient.New(bootstrap.Config)
	if err != nil {
		logger.Fatal("failed to connect to nats", zap.Error(err))
	}
	defer nc.Close()
	logger.Info("connected to NATS", zap.String("url", nc.URL()))

	callMgr := calls.NewManager()

	sipPort := os.Getenv("SIP_UDP_PORT")
	if sipPort == "" {
		sipPort = "5060"
	}

	bindAddr := os.Getenv("SIP_BIND_ADDR")
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}

	// HTTP health server
	httpPort := os.Getenv("SIP_SERVICE_PORT")
	if httpPort == "" {
		httpPort = "5070"
	}
	httpAddr := ":" + httpPort

	server := httpserver.New(httpserver.Options{
		Addr:   httpAddr,
		Logger: logger,
	})
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal("sip http server error", zap.Error(err))
		}
	}()

	// SIP UDP listener
	addr := bindAddr + ":" + sipPort
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		logger.Fatal("failed to resolve udp addr", zap.Error(err))
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Fatal("failed to listen on udp", zap.Error(err))
	}
	logger.Info("SIP UDP listener started", zap.String("addr", addr))

	buf := make([]byte, 65535)

	for {
		n, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			logger.Error("error reading udp", zap.Error(err))
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		go handleSIPMessage(authSvc, regStore, nc, callMgr, data, remote, conn, logger)
	}
}

func handleSIPMessage(
	authSvc *auth.Service,
	regStore registration.Store,
	nc *natsclient.Client,
	callMgr *calls.Manager,
	data []byte,
	remote *net.UDPAddr,
	conn *net.UDPConn,
	logger *zap.Logger,
) {
	msg, err := parser.ParseMessage(data, nil)
	if err != nil {
		logger.Warn("failed to parse sip message", zap.Error(err))
		return
	}

	request, ok := msg.(sip.Request)
	if !ok {
		logger.Warn("received non-request sip message")
		return
	}

	method := request.Method()
	logger.Info("received SIP request",
		zap.String("method", string(method)),
		zap.String("from", remote.String()),
	)

	switch method {
	case sip.REGISTER:
		handleRegister(authSvc, regStore, request, remote, conn, logger)
	case sip.INVITE:
		handleInvite(nc, callMgr, request, remote, conn, logger)
	case sip.BYE:
		handleBye(nc, callMgr, request, remote, conn, logger)
	case sip.CANCEL:
		handleCancel(nc, callMgr, request, remote, conn, logger)
	default:
		resp := sip.NewResponseFromRequest("", request, 501, "Not Implemented", "")
		sendResponse(resp, remote, conn, logger)
	}
}

func handleRegister(
	authSvc *auth.Service,
	regStore registration.Store,
	req sip.Request,
	remote *net.UDPAddr,
	conn *net.UDPConn,
	logger *zap.Logger,
) {
	logger.Info("handling SIP REGISTER", zap.String("from", remote.String()))

	// Extract username from To header
	toHeader, ok := req.To()
	if !ok || toHeader.Address == nil || toHeader.Address.User() == nil {
		logger.Warn("REGISTER missing To username")
		resp := sip.NewResponseFromRequest("", req, 400, "Bad Request", "")
		sendResponse(resp, remote, conn, logger)
		return
	}
	username := toHeader.Address.User().String()

	// Extract Contact header for contact URI
	contacts := req.GetHeaders("Contact")
	var contactURI string
	if len(contacts) > 0 {
		if ch, ok := contacts[0].(*sip.ContactHeader); ok && ch.Address != nil {
			contactURI = ch.Address.String()
		}
	}
	if contactURI == "" {
		logger.Warn("REGISTER missing Contact header")
		resp := sip.NewResponseFromRequest("", req, 400, "Bad Request", "")
		sendResponse(resp, remote, conn, logger)
		return
	}

	// Extract credentials from custom headers (temporary approach)
	var authUsername, authPassword string
	if h := req.GetHeader("X-Auth-Username"); h != nil {
		authUsername = h.Value()
	}
	if h := req.GetHeader("X-Auth-Password"); h != nil {
		authPassword = h.Value()
	}

	if authUsername == "" || authPassword == "" {
		logger.Warn("REGISTER missing auth headers")
		resp := sip.NewResponseFromRequest("", req, 401, "Unauthorized", "")
		sendResponse(resp, remote, conn, logger)
		return
	}

	// Authenticate against sip_accounts
	acc, err := authSvc.Authenticate(authUsername, authPassword)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			logger.Warn("REGISTER invalid credentials", zap.String("username", authUsername))
			resp := sip.NewResponseFromRequest("", req, 401, "Unauthorized", "")
			sendResponse(resp, remote, conn, logger)
			return
		}
		logger.Error("REGISTER auth error", zap.Error(err))
		resp := sip.NewResponseFromRequest("", req, 500, "Internal Server Error", "")
		sendResponse(resp, remote, conn, logger)
		return
	}

	// For now, tenant ID from account
	tenantID := acc.TenantID

	// Save registration binding
	expires := time.Hour
	ctx := context.Background()
	if err := regStore.SaveBinding(ctx, tenantID, username, contactURI, remote.String(), expires); err != nil {
		logger.Error("failed to save registration binding", zap.Error(err))
		resp := sip.NewResponseFromRequest("", req, 500, "Internal Server Error", "")
		sendResponse(resp, remote, conn, logger)
		return
	}

	resp := sip.NewResponseFromRequest("", req, 200, "OK", "")
	sendResponse(resp, remote, conn, logger)
}

func handleInvite(
	nc *natsclient.Client,
	callMgr *calls.Manager,
	req sip.Request,
	remote *net.UDPAddr,
	conn *net.UDPConn,
	logger *zap.Logger,
) {
	logger.Info("handling SIP INVITE", zap.String("from", remote.String()))

	// Create call ID
	callID := uuid.New().String()

	// Extract from/to usernames
	var fromUser, toUser string

	if fromHeader, ok := req.From(); ok && fromHeader.Address != nil && fromHeader.Address.User() != nil {
		fromUser = fromHeader.Address.User().String()
	}
	if toHeader, ok := req.To(); ok && toHeader.Address != nil && toHeader.Address.User() != nil {
		toUser = toHeader.Address.User().String()
	}

	if fromUser == "" || toUser == "" {
		logger.Warn("INVITE missing from/to username")
		resp := sip.NewResponseFromRequest("", req, 400, "Bad Request", "")
		sendResponse(resp, remote, conn, logger)
		return
	}

	tenantID := "tenant-default"

	now := time.Now()

	sess := &calls.CallSession{
		CallID:    callID,
		TenantID:  tenantID,
		FromUser:  fromUser,
		ToUser:    toUser,
		Direction: calls.DirectionInbound,
		CreatedAt: now,
		UpdatedAt: now,
	}
	callMgr.Create(sess)

	// Publish call started event
	startEvt := calls.CallStartedEvent{
		CallID:    callID,
		TenantID:  tenantID,
		FromUser:  fromUser,
		ToUser:    toUser,
		Direction: string(sess.Direction),
		StartedAt: now,
	}
	if err := nc.PublishJSON("calls.started", startEvt); err != nil {
		logger.Error("failed to publish calls.started", zap.Error(err))
	}

	// SIP responses: Trying then Busy (for now)
	trying := sip.NewResponseFromRequest("", req, 100, "Trying", "")
	sendResponse(trying, remote, conn, logger)

	busy := sip.NewResponseFromRequest("", req, 486, "Busy Here", "")
	sendResponse(busy, remote, conn, logger)
}

func sendResponse(resp sip.Response, remote *net.UDPAddr, conn *net.UDPConn, logger *zap.Logger) {
	data := []byte(resp.String())
	_, err := conn.WriteToUDP(data, remote)
	if err != nil {
		logger.Error("failed to send sip response", zap.Error(err))
	}
}

func handleBye(
	nc *natsclient.Client,
	callMgr *calls.Manager,
	req sip.Request,
	remote *net.UDPAddr,
	conn *net.UDPConn,
	logger *zap.Logger,
) {
	logger.Info("handling SIP BYE", zap.String("from", remote.String()))

	callID := extractCallID(req)
	if callID == "" {
		logger.Warn("BYE without Call-ID")
		resp := sip.NewResponseFromRequest("", req, 400, "Bad Request", "")
		sendResponse(resp, remote, conn, logger)
		return
	}

	callMgr.Finish(callID)

	// Publish call ended event
	tenantID := "tenant-default"
	endEvt := calls.CallEndedEvent{
		CallID:   callID,
		TenantID: tenantID,
		EndedAt:  time.Now(),
		Reason:   "BYE",
	}
	if err := nc.PublishJSON("calls.ended", endEvt); err != nil {
		logger.Error("failed to publish calls.ended", zap.Error(err))
	}

	// Respond 200 OK
	resp := sip.NewResponseFromRequest("", req, 200, "OK", "")
	sendResponse(resp, remote, conn, logger)

	// Optional: delete session
	callMgr.Delete(callID)
}

func handleCancel(
	nc *natsclient.Client,
	callMgr *calls.Manager,
	req sip.Request,
	remote *net.UDPAddr,
	conn *net.UDPConn,
	logger *zap.Logger,
) {
	logger.Info("handling SIP CANCEL", zap.String("from", remote.String()))

	callID := extractCallID(req)
	if callID == "" {
		logger.Warn("CANCEL without Call-ID")
		resp := sip.NewResponseFromRequest("", req, 400, "Bad Request", "")
		sendResponse(resp, remote, conn, logger)
		return
	}

	callMgr.Finish(callID)

	endEvt := calls.CallEndedEvent{
		CallID:   callID,
		TenantID: "tenant-default",
		EndedAt:  time.Now(),
		Reason:   "CANCEL",
	}
	if err := nc.PublishJSON("calls.ended", endEvt); err != nil {
		logger.Error("failed to publish calls.ended", zap.Error(err))
	}

	resp := sip.NewResponseFromRequest("", req, 200, "OK", "")
	sendResponse(resp, remote, conn, logger)

	callMgr.Delete(callID)
}

func extractCallID(req sip.Request) string {
	if callIDHeaders := req.GetHeaders("Call-ID"); len(callIDHeaders) > 0 {
		if h, ok := callIDHeaders[0].(*sip.CallID); ok {
			return h.String()
		}
	}
	return ""
}
