package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor представляет интерцептор для логирования
type LoggingInterceptor struct {
	logger logger
}

// NewLoggingInterceptor создает новый интерцептор логирования
func NewLoggingInterceptor(logger logger) *LoggingInterceptor {
	return &LoggingInterceptor{
		logger: logger,
	}
}

// UnaryServerInterceptor возвращает унарный серверный интерцептор для логирования
func (l *LoggingInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		l.logger.Infof("gRPC request started: method=%s", info.FullMethod)

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		code := codes.OK
		if err != nil {
			if st, ok := status.FromError(err); ok {
				code = st.Code()
			} else {
				code = codes.Unknown
			}
		}

		if err != nil {
			l.logger.Errorf("gRPC request completed: method=%s, code=%s, duration=%v, error=%v",
				info.FullMethod, code.String(), duration, err)
		} else {
			l.logger.Infof("gRPC request completed: method=%s, code=%s, duration=%v",
				info.FullMethod, code.String(), duration)
		}

		return resp, err
	}
}
