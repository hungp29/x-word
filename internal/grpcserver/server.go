package grpcserver

import (
	"context"
	"log/slog"
	"time"

	wordv1 "github.com/hungp29/x-proto/gen/go/word/v1"
	"github.com/hungp29/x-word/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WordServiceServer implements wordv1.WordServiceServer using the application WordService.
type WordServiceServer struct {
	wordv1.UnimplementedWordServiceServer
	svc    *service.WordService
	logger *slog.Logger
}

// NewWordServiceServer returns a gRPC server that delegates to the given WordService.
func NewWordServiceServer(svc *service.WordService, logger *slog.Logger) *WordServiceServer {
	return &WordServiceServer{svc: svc, logger: logger}
}

// GetWord returns a single word entry.
func (s *WordServiceServer) GetWord(ctx context.Context, req *wordv1.GetWordRequest) (*wordv1.GetWordResponse, error) {
	if req == nil || req.Word == "" {
		return nil, status.Error(codes.InvalidArgument, "word is required")
	}
	dict := serviceDictFromProto(req.Dict)
	w, err := s.svc.GetWord(req.Word, dict)
	if err != nil {
		s.logger.Error("GetWord failed", "word", req.Word, "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &wordv1.GetWordResponse{Word: wordToProto(w)}, nil
}

// GetWords returns multiple word entries.
func (s *WordServiceServer) GetWords(ctx context.Context, req *wordv1.GetWordsRequest) (*wordv1.GetWordsResponse, error) {
	if req == nil || len(req.Words) == 0 {
		return &wordv1.GetWordsResponse{Words: nil}, nil
	}
	dict := serviceDictFromProto(req.Dict)
	words, err := s.svc.GetWords(req.Words, dict)
	if err != nil {
		s.logger.Error("GetWords failed", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := make([]*wordv1.Word, len(words))
	for i, w := range words {
		out[i] = wordToProto(w)
	}
	return &wordv1.GetWordsResponse{Words: out}, nil
}

// Ping returns service info (replaces HTTP GET /).
func (s *WordServiceServer) Ping(ctx context.Context, req *wordv1.PingRequest) (*wordv1.PingResponse, error) {
	return &wordv1.PingResponse{
		Service: "x-word",
		Message: "hello from x-word",
		Time:    time.Now().Format(time.RFC3339),
	}, nil
}

// Health returns health status (replaces HTTP GET /health).
func (s *WordServiceServer) Health(ctx context.Context, req *wordv1.HealthRequest) (*wordv1.HealthResponse, error) {
	return &wordv1.HealthResponse{Status: "ok"}, nil
}
