package main

import (
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
	
	// TODO: set server and set querier as param
	fmt.Println("Hello, Posterr!")
}
