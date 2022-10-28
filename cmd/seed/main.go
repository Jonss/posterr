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
		users := []string{"rhaenyra", "viserys", "daemon", "aegon"}
		execSeed(ctx, cfg, querier, users)
	}
}

func execSeed(
	ctx context.Context,
	cfg config.Config,
	querier *db.Queries,
	users []string,
) {
	for _, u := range users {
		if _, err := querier.SeedUser(ctx, u); err != nil {
			log.Fatalf("seed users failed: error=(%v)", err)
		}
	}
}
