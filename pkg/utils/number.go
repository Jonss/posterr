package utils

import "database/sql"

func Int64PtrToNullInt64(n *int64) sql.NullInt64 {
	if n == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *n, Valid: true}
}

func NullInt64ToInt64Ptr(n sql.NullInt64) *int64 {
	if !n.Valid {
		return nil
	}
	return &n.Int64
}

func NullInt64ToInt64(n sql.NullInt64) int64 {
	if !n.Valid {
		return 0
	}
	return n.Int64
}

func Int64ToNullInt64(n int64) sql.NullInt64 {
	return Int64PtrToNullInt64(&n)
}

func Int64ToPtr(n int64) *int64 {
	return &n
}
