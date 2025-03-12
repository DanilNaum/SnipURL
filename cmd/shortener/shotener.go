package main

import (
	"context"
	"errors"
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
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	// создаём предустановленный регистратор zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar := logger.Sugar()
	sugar.Info("App is running...")

	err = run(sugar)

	if err != nil && !errors.Is(err, context.Canceled) {
		sugar.Fatal("App fail with error %s", err)
	}

	sugar.Info("App is gracefully shutdown")
	os.Exit(0)
}

func run(log *zap.SugaredLogger) error {

	conf := config.NewConfig(log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	services, ctx := errgroup.WithContext(ctx)

	storage := memory.NewStorage()

	hash := hash.NewHasher(8)

	urlSnipperService := urlsnipper.NewURLSnipperService(storage, hash)

	mux := chi.NewRouter()

	controller, err := rest.NewController(mux, conf.ServerConfig(), urlSnipperService, log)

	if err != nil {
		return err
	}

	httpServer := httpserver.NewHTTPServer(controller, httpserver.WithAddr(conf.ServerConfig().HTTPServerHost()))

	services.Go(func() error {
		err := <-httpServer.Notify()
		return err
	})

	go func() {
		<-ctx.Done()
		err := httpServer.Shutdown()
		if err != nil {
			log.Errorf("shutdown error: %s", err)
		}
	}()

	err = services.Wait()

	if err == nil || errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}
