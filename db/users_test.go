package db

import (
	"context"
	"database/sql"
	"testing"
)

func TestSeedUsers(t *testing.T) {
	querier, tearDown := newDbTestSetup(t)
	defer tearDown()

	ctx := context.Background()

	users := []string{
		"fferdinand",
		"fjosef",
		"sissi",
		"maximilian",
	}
	for _, u := range users {
		if _, err := querier.SeedUser(ctx, u); err != nil {
			t.Fatalf("error seeding users: error=(%v)", err)
		}
	}

	testCases := []struct {
		name     string
		username string
		wantErr  error
	}{
		{
			name:     "should find user fferdinand",
			username: "fferdinand",
		},
		{
			name:     "should find user fjosef",
			username: "fjosef",
		},
		{
			name:     "should find user sissi",
			username: "sissi",
		},
		{
			name:     "should find user maximilian",
			username: "maximilian",
		},
		{
			name:     "should not find user connor_mcleod",
			username: "connor_mcleod",
			wantErr:  sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := "select id, username from users where username = $1"
			row := querier.db.QueryRowContext(ctx, query, tc.username)
			var got User
			err := row.Scan(&got.ID, &got.Username)

			if err != nil && err != tc.wantErr {
				t.Fatalf("err = %v, want %v", tc.wantErr, err)
			}

			if err == nil {
				if got.Username != tc.username {
					t.Fatalf("username = %v, want %v", got.Username, tc.username)
				}
			}
		})
	}
}
