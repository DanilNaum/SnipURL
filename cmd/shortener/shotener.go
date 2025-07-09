package main

import (
	"context"
	"errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/DanilNaum/SnipURL/internal/app/config"
	"github.com/DanilNaum/SnipURL/internal/app/repository/url/memory"
	"github.com/DanilNaum/SnipURL/internal/app/repository/url/psql"
	"github.com/DanilNaum/SnipURL/internal/app/service/urlsnipper"
	rest "github.com/DanilNaum/SnipURL/internal/app/transport/rest"
	"github.com/DanilNaum/SnipURL/pkg/cookie"
	"github.com/DanilNaum/SnipURL/pkg/migration"
	"github.com/DanilNaum/SnipURL/pkg/pg"
	"github.com/DanilNaum/SnipURL/pkg/utils/dumper"
	"github.com/DanilNaum/SnipURL/pkg/utils/hash"
	"github.com/DanilNaum/SnipURL/pkg/utils/httpserver"
	"github.com/go-chi/chi/v5"

	_ "net/http/pprof"

	urlstorage "github.com/DanilNaum/SnipURL/internal/app/repository/url"
	deleteurl "github.com/DanilNaum/SnipURL/internal/app/service/delete"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("error create logger %s", err.Error())
	}

	sugarLogger := logger.Sugar()

	sugarLogger.Info("Build version: ", buildVersion, "\nBuild date: ", buildDate, "\nBuild commit: ", buildCommit)

	err = run(sugarLogger)

	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Sync()
		sugarLogger.Fatalf("App fail with error %s", err.Error())
	}

	sugarLogger.Info("App is gracefully shutdown")
	logger.Sync()

}

func run(log *zap.SugaredLogger) error {

	conf := config.NewConfig(log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	defer cancel()

	dump, err := dumper.NewDumper(conf.DumpConfig().GetPath(), log)
	if err != nil {
		return err
	}

	var urlStorage urlstorage.URLStorage

	if conf.DBConfig().GetDSN() != "" {
		migrator := migration.NewMigrator(conf.DBConfig().GetDSN(), migration.WithRelativePath("migrations"))
		err = migrator.Migrate()
		if err != nil {
			return err
		}

		pgConf := pg.NewConnConfigFromDsnString(conf.DBConfig().GetDSN())

		pgConn := pg.NewConnection(ctx, pgConf, log)
		if pgConn == nil {
			return errors.New("pg connection is nil")
		}
		defer pgConn.Close()
		urlStorage = psql.NewStorage(pgConn)
	} else {
		storage := memory.NewStorage()

		err = storage.RestoreStorage(dump)
		if err != nil {
			return err
		}
		urlStorage = storage
		defer dump.Close()
	}

	services, ctx := errgroup.WithContext(ctx)

	hash := hash.NewHasher(8)

	deleteService := deleteurl.NewDeleteService(ctx, urlStorage)
	urlSnipperService := urlsnipper.NewURLSnipperService(urlStorage, hash, dump, deleteService, log)

	mux := chi.NewRouter()

	cookieManager := cookie.NewCookieManager([]byte(conf.CookieConfig().GetSecret()), cookie.WithName("user"))

	controller, err := rest.NewController(mux, conf.ServerConfig(), urlSnipperService, urlStorage, cookieManager, log)

	if err != nil {
		return err
	}

	httpServer := httpserver.NewHTTPServer(controller, httpserver.WithAddr(conf.ServerConfig().HTTPServerHost()), httpserver.WithTLS(conf.ServerConfig().GetEnableHTTPs()))

	services.Go(func() error {
		err = <-httpServer.Notify()
		return err
	})

	go func() {
		<-ctx.Done()
		err = httpServer.Shutdown()
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
