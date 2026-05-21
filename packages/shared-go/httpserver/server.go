package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	logger *zap.Logger
	server *http.Server
}

type Options struct {
	Addr   string
	Logger *zap.Logger
}

func New(opts Options) *Server {
	if opts.Addr == "" {
		opts.Addr = ":8080"
	}

	if opts.Logger == nil {
		panic("httpserver: logger is required")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: router,
	}

	return &Server{
		engine: router,
		logger: opts.Logger,
		server: s,
	}
}

// Engine exposes the Gin engine so services can register additional routes.
func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Start() error {
	go func() {
		s.logger.Info("http server starting", zap.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("http server error", zap.Error(err))
		}
	}()

	// graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	s.logger.Info("received shutdown signal", zap.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("http server shutdown error", zap.Error(err))
		return err
	}

	s.logger.Info("http server stopped gracefully")
	return nil
}
