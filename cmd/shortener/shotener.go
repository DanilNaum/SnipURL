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
	"github.com/DanilNaum/SnipURL/pkg/utils/dumper"
	"github.com/DanilNaum/SnipURL/pkg/utils/hash"
	"github.com/DanilNaum/SnipURL/pkg/utils/httpserver"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}

	defer logger.Sync()

	sugarLogger := logger.Sugar()

	sugarLogger.Info("App is running...")

	err = run(sugarLogger)

	if err != nil && !errors.Is(err, context.Canceled) {
		sugarLogger.Fatalf("App fail with error %s", err.Error())
	}

	sugarLogger.Info("App is gracefully shutdown")
	os.Exit(0)
}

func run(log *zap.SugaredLogger) error {

	conf := config.NewConfig(log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	defer cancel()

	services, ctx := errgroup.WithContext(ctx)

	storage := memory.NewStorage()

	hash := hash.NewHasher(8)

	dump, err := dumper.NewDumper(conf.DumpConfig().GetPath(), log)
	if err != nil {
		return err
	}

	defer dump.Close()

	urlSnipperService := urlsnipper.NewURLSnipperService(storage, hash, dump, log)

	err = urlSnipperService.RestoreStorage()
	if err != nil {
		return err
	}

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
