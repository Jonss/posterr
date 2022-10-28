package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Jonss/posterr/api/httpapi"
	"github.com/Jonss/posterr/config"
	"github.com/Jonss/posterr/db"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: error=(%v)", err)
	}

	dbConn, err := db.NewConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: error=(%v)", err)
	}
	db.New(dbConn)

	if cfg.ShouldMigrate {
		if err := db.Migrate(dbConn, cfg.DBName, cfg.MigrationPath); err != nil {
			log.Fatalf("key 'SHOULD_MIGRATE' is (%t) but migration failed: error=(%v)", cfg.ShouldMigrate, err)
		}
	}

	router := mux.NewRouter()

	services := httpapi.Services{}
	httpServer := httpapi.NewHttpServer(
		router,
		cfg,
		services,
	)
	httpServer.Start()

	addr := "0.0.0.0:"+cfg.Port
	server := &http.Server{
		Handler: router,
		Addr:   addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Posterr Server started!")
	log.Fatal(server.ListenAndServe())
}
