package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"

	"github.com/DanilNaum/SnipURL/internal/app/repository/url/memory"
	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	rest "github.com/DanilNaum/SnipURL/internal/app/transport/rest"
	"github.com/DanilNaum/SnipURL/pkg/utils/hash"
	"github.com/DanilNaum/SnipURL/pkg/utils/httpserver"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := run()

	if err != nil && !errors.Is(err, context.Canceled) {
		// TODO: Log
		os.Exit(1)
	}

	//TODO: Log
	os.Exit(0)
}

func run() error {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	services, _ := errgroup.WithContext(ctx)

	storage := memory.NewStorage()

	hash := hash.NewHasher(8)

	urlSnipperService := urlsnipper.NewURLSnipperService(storage, hash)

	mux := http.NewServeMux()

	controller := rest.NewController(mux, urlSnipperService)

	httpServer := httpserver.NewHTTPServer(controller)

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

	err := services.Wait()

	if err == nil || errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}
