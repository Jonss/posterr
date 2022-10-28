package db

import (
	"context"
	"testing"
	"time"

	"github.com/jonss/testcontainers-go-wrapper/pg"
)

/* newDbTestSetup does:
- configures a container with postgres,
- creates a connection
- run migration
- returns a New Queries to
*/
func newDbTestSetup(t *testing.T) (*Queries, func()) {
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

// parsedDate parses a date in the format yyyy-MM-dd
func parsedDate(dateStr string, t *testing.T) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t.Fatalf("error creating users. error=(%v)", err)
	}
	return date
}

// parsedDatePtr parses a date in the format yyyy-MM-dd as pointer
func parsedDatePtr(dateStr string, t *testing.T) *time.Time {
	d := parsedDate(dateStr, t)
	return &d
}
