package interceptors

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// TrustedSubnetInterceptor представляет интерцептор для проверки доверенной подсети
type TrustedSubnetInterceptor struct {
	trustedSubnet *net.IPNet
	logger        logger
}

// NewTrustedSubnetInterceptor создает новый интерцептор проверки доверенной подсети
func NewTrustedSubnetInterceptor(trustedSubnetCIDR string, logger logger) (*TrustedSubnetInterceptor, error) {
	if trustedSubnetCIDR == "" {
		return &TrustedSubnetInterceptor{
			trustedSubnet: nil,
			logger:        logger,
		}, nil
	}

	_, subnet, err := net.ParseCIDR(trustedSubnetCIDR)
	if err != nil {
		return nil, err
	}

	return &TrustedSubnetInterceptor{
		trustedSubnet: subnet,
		logger:        logger,
	}, nil
}

// UnaryServerInterceptor возвращает унарный серверный интерцептор для проверки доверенной подсети
func (t *TrustedSubnetInterceptor) UnaryServerInterceptor(protectedMethods map[string]bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if !protectedMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		if t.trustedSubnet == nil {
			return nil, status.Errorf(codes.PermissionDenied, "access denied")
		}

		clientIP, err := t.getClientIP(ctx)
		if err != nil {
			t.logger.Errorf("Failed to get client IP for method %s: %v", info.FullMethod, err)
			return nil, status.Errorf(codes.PermissionDenied, "access denied")
		}

		if !t.trustedSubnet.Contains(clientIP) {
			t.logger.Errorf("Access denied for IP %s to method %s: not in trusted subnet %s",
				clientIP.String(), info.FullMethod, t.trustedSubnet.String())
			return nil, status.Errorf(codes.PermissionDenied, "access denied")
		}

		t.logger.Infof("Access granted for IP %s to method %s", clientIP.String(), info.FullMethod)
		return handler(ctx, req)
	}
}

func (t *TrustedSubnetInterceptor) getClientIP(ctx context.Context) (net.IP, error) {
	// Сначала пытаемся получить IP из метаданных (X-Real-IP, X-Forwarded-For)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// Проверяем X-Real-IP
		if realIPs := md.Get("x-real-ip"); len(realIPs) > 0 && realIPs[0] != "" {
			ip := net.ParseIP(realIPs[0])
			if ip != nil {
				return ip, nil
			}
		}

		// Проверяем X-Forwarded-For
		if forwardedIPs := md.Get("x-forwarded-for"); len(forwardedIPs) > 0 && forwardedIPs[0] != "" {
			// X-Forwarded-For может содержать несколько IP, разделенных запятыми
			ips := strings.Split(forwardedIPs[0], ",")
			if len(ips) > 0 {
				ip := net.ParseIP(strings.TrimSpace(ips[0]))
				if ip != nil {
					return ip, nil
				}
			}
		}
	}

	// Если не удалось получить IP из метаданных, используем peer info
	if p, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := p.Addr.(*net.TCPAddr); ok {
			return tcpAddr.IP, nil
		}
	}

	return nil, status.Errorf(codes.Internal, "failed to get client IP")
}
