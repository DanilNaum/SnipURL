package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DanilNaum/SnipURL/internal/app/config"
	"github.com/DanilNaum/SnipURL/internal/app/repository/url/memory"
	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	rest "github.com/DanilNaum/SnipURL/internal/app/transport/rest"
	"github.com/DanilNaum/SnipURL/pkg/utils/hash"
	"github.com/DanilNaum/SnipURL/pkg/utils/httpserver"
	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := log.New(os.Stdout, "URL Snipper: ", log.LstdFlags)
	logger.Println("App is running...")
	err := run(logger)

	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Printf("App fail with error %w", err)
		os.Exit(1)
	}

	logger.Println("App is gracefully shutdown")
	os.Exit(0)
}

func run(log *log.Logger) error {

	conf := config.GetConfig(log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	services, ctx := errgroup.WithContext(ctx)

	storage := memory.NewStorage()

	hash := hash.NewHasher(8)

	urlSnipperService := urlsnipper.NewURLSnipperService(storage, hash)

	mux := chi.NewRouter()

	controller, err := rest.NewController(mux, conf.ServerConfig(), urlSnipperService)

	if err != nil {
		return err
	}

	httpServer := httpserver.NewHTTPServer(controller, httpserver.WithAddr(conf.ServerConfig().HttpServerHost()))

	services.Go(func() error {
		err := <-httpServer.Notify()
		return err
	})

	go func() {
		<-ctx.Done()
		err := httpServer.Shutdown()
		if err != nil {
			log.Printf("shutdown error: %s", err)
		}
	}()

	err = services.Wait()

	if err == nil || errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}
