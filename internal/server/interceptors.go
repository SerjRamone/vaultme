package server

import (
	context "context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) requestLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			s.log.Error("processing request error", zap.Error(err))
			return nil, fmt.Errorf("processing request error: %w", err)
		} else {
			md, ok := metadata.FromIncomingContext(ctx)
			if ok {
				s.log.Info("request",
					zap.String("method", info.FullMethod),
					zap.Any("headers", md),
					zap.Duration("duration", duration),
					zap.Any("body", req),
				)
			}
		}
		return resp, nil
	}
}
