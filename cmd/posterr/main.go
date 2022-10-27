package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Jonss/posterr/config"
	"github.com/Jonss/posterr/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: error=(%v)", err)
	}

	ctx := context.Background()
	dbConn, err := db.NewConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: error=(%v)", err)
	}
	querier := db.New(dbConn)
	handleDatabase(ctx, cfg, dbConn, querier)
	
	// TODO: set server and set querier as param
	fmt.Println("Hello, Posterr!")
}

func handleDatabase(
	ctx context.Context,
	cfg config.Config,
	dbConn *sql.DB,
	querier *db.Queries,
) {
	if cfg.ShouldMigrate {
		if err := db.Migrate(dbConn, cfg.DBName, cfg.MigrationPath); err != nil {
			log.Fatalf("key 'SHOULD_MIGRATE' is (%t) but migration failed: error=(%v)", cfg.ShouldMigrate, err)
		}
	}

	if cfg.IsLocal() {
		if err := querier.SeedUsers(ctx); err != nil {
			log.Fatalf("seed users failed: error=(%v)", err)
		}
	}
}