package strings

import "database/sql"

func NullStrToPointer(str sql.NullString) *string {
	if !str.Valid {
		return nil
	}
	return &str.String
}

func StrToNullStr(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}
