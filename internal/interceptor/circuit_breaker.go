package interceptor

import (
	"context"
	"user_service/internal/logger"

	"github.com/sony/gobreaker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerInterceptor struct {
	cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker(cb *gobreaker.CircuitBreaker) *CircuitBreakerInterceptor {
	return &CircuitBreakerInterceptor{cb: cb}
}

func (c *CircuitBreakerInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := c.cb.Execute(func() (interface{}, error) {
		return handler(ctx, req)
	})
	if err != nil {
		if err == gobreaker.ErrOpenState {
			logger.Error("service unavailable", zap.Error(err))
			return nil, status.Error(codes.Unavailable, "service unavailable")
		}
		return nil, err
	}
	return res, nil
}
