package db

import (
	"context"
	"database/sql"
	"testing"
)

func TestSeedUsers(t *testing.T) {
	querier, tearDown := NewDbTestSetup(t)
	defer tearDown()

	ctx := context.Background()
	if err := querier.SeedUsers(ctx); err != nil {
		t.Fatalf("error seeding users: error=(%v)", err)
	}

	testCases := []struct{
		name string
		username string
		wantID int
		wantErr error
	} {
		{
			name: "should find user fferdinand",
			username: "fferdinand",
			wantID: 1,
		},
		{
			name: "should find user fjosef",
			username: "fjosef",
			wantID: 2,
		},
		{
			name: "should find user sissi",
			username: "sissi",
			wantID: 3,
		},
		{
			name: "should find user maximilian",
			username: "maximilian",
			wantID: 4,
		},
		{
			name: "should not find user connor_mcleod",
			username: "connor_mcleod",
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := "select id, username from users where username = $1"
			row := querier.db.QueryRowContext(ctx, query, tc.username)
			var user User
			err := row.Scan(&user.ID, &user.Username)
		
			if err != nil && err != tc.wantErr {
				t.Fatalf("err = %v, want %v", tc.wantErr, err)
			}

			if err == nil {
				if user.ID != int64(tc.wantID) {
					t.Fatalf("err = %v, want %v", user.ID, tc.wantID)
				}
	
				if user.Username != tc.username {
					t.Fatalf("username = %v, want %v", user.Username, tc.username)
				}
			}	
		})
	}
}