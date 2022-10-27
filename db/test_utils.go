package db

import (
	"context"
	"testing"

	"github.com/jonss/testcontainers-go-wrapper/pg"
)

/* NewDbTestSetup does:
- configures a container with postgres,
- creates a connection
- run migration
- returns a New Queries to
*/
func NewDbTestSetup(t *testing.T) (*Queries, func()) {
	cfg := pg.PostgresCfg{
		ImageName: "postgres:15-alpine",
		Password:  "a_secret_password",
		UserName:  "test",
		DbName:    "posterr_test",
	}
	
	pgInfo, err := pg.Container(context.Background(), cfg)
	if err != nil {
		t.Fatalf("error creating pgContainer. error=(%v)", err)
	}

	dbConn, err := NewConnection(pgInfo.DbURL)
	if err != nil {
		t.Fatalf("error connecting db. error=(%v)", err)
	}

	err = Migrate(dbConn,
		cfg.DbName,
		"migrations",
	)

	if err != nil {
		t.Fatalf("error connecting db. error=(%v)", err)
	}

	return New(dbConn), pgInfo.TearDown
}