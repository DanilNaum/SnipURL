package grpc

import (
	"context"
	"net"

	"github.com/DanilNaum/SnipURL/internal/app/transport/grpc/interceptors"
	"github.com/DanilNaum/SnipURL/pkg/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type cookieManager interface {
	GetFromMetadata(md metadata.MD) (string, error)
	SetToMetadata(userID string) metadata.MD
}

type logger interface {
	Infof(format string, v ...any)
	Errorf(format string, v ...any)
}

// Controller представляет gRPC контроллер
type Controller struct {
	server        *grpc.Server
	snipURLServer *Server
	logger        logger
}

// NewController создает новый gRPC контроллер с настроенными интерцепторами
func NewController(
	service service,
	internalService internalService,
	psqlStoragePinger psqlStoragePinger,
	conf config,
	cookieManager cookieManager,
	logger logger,
	trustedSubnetCIDR string,
) (*Controller, error) {
	authInterceptor := interceptors.NewAuthInterceptor(cookieManager, logger)
	loggingInterceptor := interceptors.NewLoggingInterceptor(logger)

	trustedSubnetInterceptor, err := interceptors.NewTrustedSubnetInterceptor(trustedSubnetCIDR, logger)
	if err != nil {
		return nil, err
	}

	protectedAuthMethods := map[string]bool{
		"/snipurl.SnipURLService/GetUserURLs":    true,
		"/snipurl.SnipURLService/DeleteUserURLs": true,
	}

	protectedSubnetMethods := map[string]bool{
		"/snipurl.SnipURLService/GetStats": true,
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor.UnaryServerInterceptor(),
			authInterceptor.UnaryServerInterceptor(),
			interceptors.RequireAuthInterceptor(protectedAuthMethods, logger),
			trustedSubnetInterceptor.UnaryServerInterceptor(protectedSubnetMethods),
		),
	)

	snipURLServer, err := NewServer(service, internalService, psqlStoragePinger, conf)
	if err != nil {
		return nil, err
	}

	protobuf.RegisterSnipURLServiceServer(server, snipURLServer)

	return &Controller{
		server:        server,
		snipURLServer: snipURLServer,
		logger:        logger,
	}, nil
}

func (c *Controller) serve(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		c.logger.Errorf("Failed to listen on %s: %v", address, err)
		return err
	}
	reflection.Register(c.server)
	c.logger.Infof("Starting gRPC server on %s", address)

	if err := c.server.Serve(listener); err != nil {
		c.logger.Errorf("Failed to serve gRPC server: %v", err)
		return err
	}

	return nil
}

func (c *Controller) stop() {
	c.logger.Infof("Stopping gRPC server")
	c.server.GracefulStop()
}

// ServeWithCtx запускает gRPC сервер на указанном адресе, который завершиться при отмене контекста

func (c *Controller) ServeWithCtx(ctx context.Context, address string) chan error {
	errChan := make(chan error)

	go func() {
		errChan <- c.serve(address)
		close(errChan)
	}()

	go func() {
		<-ctx.Done()
		c.stop()
	}()
	return errChan
}
