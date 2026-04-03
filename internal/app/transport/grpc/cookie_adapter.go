package grpc

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type CookieManager struct{}

// NewCookieManager создает новый адаптер для cookie manager
func NewCookieManager() *CookieManager {
	return &CookieManager{}
}

// GetFromMetadata извлекает userID из gRPC метаданных
func (c *CookieManager) GetFromMetadata(md metadata.MD) (string, error) {

	cookies := md.Get("cookie")
	if len(cookies) == 0 {
		return "", errors.New("no cookies found in metadata")
	}

	for _, cookie := range cookies {

		if len(cookie) > 8 && cookie[:8] == "user_id=" {
			userID := cookie[8:]
			if userID != "" {
				return userID, nil
			}
		}
	}

	return "", errors.New("user_id not found in cookies")
}

// SetToMetadata создает метаданные с cookie для userID
func (c *CookieManager) SetToMetadata(userID string) metadata.MD {

	cookieValue := fmt.Sprintf("user_id=%s; Path=/; HttpOnly", userID)

	return metadata.Pairs("set-cookie", cookieValue)
}

// GetFromGRPCContext извлекает userID из gRPC контекста через метаданные
func (c *CookieManager) GetFromGRPCContext(md metadata.MD) (string, error) {
	return c.GetFromMetadata(md)
}

// SetToGRPCContext устанавливает userID в gRPC контекст через метаданные
func (c *CookieManager) SetToGRPCContext(userID string) metadata.MD {
	return c.SetToMetadata(userID)
}
