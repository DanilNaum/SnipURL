package interceptors

import (
	"context"

	"github.com/DanilNaum/SnipURL/internal/app/transport/rest/middlewares"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type cookieManager interface {
	GetFromMetadata(md metadata.MD) (string, error)
	SetToMetadata(userID string) metadata.MD
}

type logger interface {
	Infof(format string, v ...any)
	Errorf(format string, v ...any)
}

// AuthInterceptor представляет интерцептор для аутентификации
type AuthInterceptor struct {
	cookieManager cookieManager
	logger        logger
}

// NewAuthInterceptor создает новый интерцептор аутентификации
func NewAuthInterceptor(cookieManager cookieManager, logger logger) *AuthInterceptor {
	return &AuthInterceptor{
		cookieManager: cookieManager,
		logger:        logger,
	}
}

var key = middlewares.Key{Key: "userID"}

// UnaryServerInterceptor возвращает унарный серверный интерцептор для аутентификации
func (a *AuthInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			a.logger.Errorf("Failed to get metadata from context for method %s", info.FullMethod)
			return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
		}

		userID, err := a.cookieManager.GetFromMetadata(md)
		if err != nil {
			a.logger.Infof("No valid user ID found in metadata for method %s, creating new user", info.FullMethod)
			// Если пользователь не найден, создаем нового
			// В реальном приложении здесь может быть логика создания нового пользователя
			userID = uuid.NewString() // Заглушка для нового пользователя
		}

		newCtx := context.WithValue(ctx, key, userID)

		resp, err := handler(newCtx, req)

		if userID != "" {
			outgoingMD := a.cookieManager.SetToMetadata(userID)
			grpc.SetHeader(newCtx, outgoingMD)
		}

		return resp, err
	}
}

// RequireAuthInterceptor создает интерцептор, который требует аутентификации для определенных методов
func RequireAuthInterceptor(protectedMethods map[string]bool, logger logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if protectedMethods[info.FullMethod] {
			userID := ctx.Value(key)
			if userID == nil || userID == "" {
				logger.Errorf("Authentication required for method %s but no user ID found", info.FullMethod)
				return nil, status.Errorf(codes.Unauthenticated, "authentication required")
			}
		}

		return handler(ctx, req)
	}
}
