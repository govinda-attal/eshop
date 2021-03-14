package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-logrusutil/logrusutil/logctx"
	hc "github.com/heptiolabs/healthcheck"
	"github.com/spf13/cobra"

	"github.com/govinda-attal/eshop/internal/eshop"
	"github.com/govinda-attal/eshop/internal/platform/store"
	"github.com/govinda-attal/eshop/internal/platform/web"
	"github.com/govinda-attal/eshop/pkg/eshop/api"
)

// serveCmd represents the migrate command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the api",
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	var (
		ctx = logctx.New(context.Background(), log)
	)
	db, err := store.Crdb(ctx, cfg.DB.Url)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Info("crdb ping success")

	srv := eshop.NewApi(ctx, eshop.NewCartService(ctx, db))
	ws, err := setupWebServer(ctx, srv)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := ws.Start(ctx); err != nil {
			log.Fatal("eshop service startup failed: ", err)
		}
	}()
	log.Info("eshop service started ...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Infof("eshop service shutdown in progress with grace period of %s ...", cfg.ShutdownGracePeriod)
	ctx, cancel := context.WithTimeout(ctx, cfg.ShutdownGracePeriod)
	defer cancel()

	<-ctx.Done()
	log.Info("eshop service shutdown ...")
	os.Exit(0)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func setupWebServer(ctx context.Context, srv api.ServerInterface) (*web.Server, error) {
	health, err := healthHandler()
	if err != nil {
		return nil, err
	}
	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}
	ws := web.NewServer(ctx, cfg.Port,
		web.ServerWithGraceShutdown(cfg.ShutdownGracePeriod),
		web.ServerWithHealthCheck("/health", health),
		web.ServerWithPrometheusHandler("/metrics"),
		web.ServerWithApiSpec(swagger),
	)
	api.RegisterHandlersWithBaseURL(ws.Echo(), srv, "/api")
	return ws, nil
}

func healthHandler() (http.HandlerFunc, error) {
	return hc.NewHandler().ReadyEndpoint, nil
}
