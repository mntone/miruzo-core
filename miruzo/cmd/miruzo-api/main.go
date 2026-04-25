package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mntone/miruzo-core/miruzo/internal/adapter/persistence"
	m "github.com/mntone/miruzo-core/miruzo/internal/api/middleware"
	"github.com/mntone/miruzo-core/miruzo/internal/app"
	"github.com/mntone/miruzo-core/miruzo/internal/server"
)

var version = "0.0.0+dev"

func main() {
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	hdl, err := persistence.OpenAppHandle(context.Background(), cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer hdl.Close()

	mux := http.NewServeMux()
	app.MountAll(mux, hdl.PersistenceProvider(), cfg, version)

	httpServer := server.NewHTTPServer(
		m.RequestLog(mux),
		cfg.Server,
	)
	log.Printf("[miruzo-api] version %s", version)
	log.Printf("[miruzo-api] listening on %s", httpServer.Addr)
	if err := server.Run(httpServer, cfg.Server.ShutdownTimeout); err != nil {
		log.Fatal(err)
	}
}
