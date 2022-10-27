package main

import (
	"context"
	"log"

	"github.com/Jonss/posterr/config"
	"github.com/Jonss/posterr/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("[seed] error loading config: error=(%v)", err)
	}

	ctx := context.Background()
	dbConn, err := db.NewConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("[seed] error connecting to database: error=(%v)", err)
	}
	querier := db.New(dbConn)
	if cfg.IsLocal() {
		execSeed(ctx, cfg, querier)
	}
}


func execSeed(
	ctx context.Context,
	cfg config.Config,
	querier *db.Queries,
) {
	if err := querier.SeedUsers(ctx); err != nil {
		log.Fatalf("seed users failed: error=(%v)", err)
	}
}