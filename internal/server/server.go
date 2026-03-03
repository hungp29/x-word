package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hupham/x-word/internal/config"
)

// Server wraps the Gin engine and config. No shared mutable request/session state.
type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	logger *slog.Logger
}

// New builds a Server with the given config and logger. Router and handlers are set up here.
func New(cfg *config.Config, logger *slog.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(requestLogger(logger))

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "x-word", "message": "hello"})
	})

	return &Server{engine: engine, cfg: cfg, logger: logger}
}

// Run starts the HTTP server and blocks until ctx is cancelled or the server fails.
func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.HTTPPort),
		Handler: s.engine,
	}
	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()
	s.logger.Info("http server listening", "port", s.cfg.HTTPPort)
	return srv.ListenAndServe()
}

// requestLogger returns a Gin middleware that logs request method, path, and status (structured).
func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		status := c.Writer.Status()
		logger.Info("request", "method", method, "path", path, "status", status)
	}
}
