package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/hungp29/x-proto/gen/go/word/v1"
	"github.com/hungp29/x-word/internal/config"
	"github.com/hungp29/x-word/internal/fetcher"
	"github.com/hungp29/x-word/internal/grpcserver"
	"github.com/hungp29/x-word/internal/parser"
	"github.com/hungp29/x-word/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server wraps the gRPC server and config. No shared mutable request/session state.
type Server struct {
	grpc   *grpc.Server
	cfg    *config.Config
	logger *slog.Logger
}

// New builds a Server with the given config and logger. Registers WordService and reflection.
func New(cfg *config.Config, logger *slog.Logger) *Server {
	grpcOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(unaryLoggingInterceptor(logger)),
	}
	s := grpc.NewServer(grpcOpts...)

	f := fetcher.NewHTTPFetcher()
	wordSvc := service.NewWordService(f, map[service.Dictionary]service.Parser{
		service.DictionaryEnglish:           parser.NewCambridgeParser(),
		service.DictionaryEnglishVietnamese: parser.NewEnglishVietnameseParser(),
	})
	wordv1.RegisterWordServiceServer(s, grpcserver.NewWordServiceServer(wordSvc, logger))
	reflection.Register(s)

	return &Server{grpc: s, cfg: cfg, logger: logger}
}

// Run starts the gRPC server and blocks until ctx is cancelled or the server fails.
func (s *Server) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		stopped := make(chan struct{})
		go func() {
			s.grpc.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-shutdownCtx.Done():
			s.grpc.Stop()
		}
	}()
	s.logger.Info("grpc server listening", "port", s.cfg.GRPCPort)
	return s.grpc.Serve(lis)
}

// unaryLoggingInterceptor logs method and status for each unary RPC.
func unaryLoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		status := "ok"
		if err != nil {
			status = "error"
		}
		logger.Info("rpc", "method", info.FullMethod, "status", status, "duration_ms", time.Since(start).Milliseconds())
		return resp, err
	}
}
