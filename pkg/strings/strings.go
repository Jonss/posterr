package strings

import (
	"database/sql"
	"net/url"
	"strconv"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

func NullStrToPointer(str sql.NullString) *string {
	if !str.Valid {
		return nil
	}
	return &str.String
}

func StrToNullStr(str string) sql.NullString {
	return sql.NullString{String: str, Valid: true}
}

func ParseStringToDate(values url.Values, key string ) (*time.Time, error) {
	value := values[key]
	if len(value) == 0 {
		return nil, nil
	}
	
	endDate, err := time.Parse(dateFormat, value[0])
	if err != nil {
		return &time.Time{}, err
	}
	return &endDate, nil
}

func ParseStringToInt(values url.Values, key string) (int, error) {
	value := values[key]
	if len(value) == 0 {
		return 0, nil
	}
	n, err := strconv.Atoi(value[0])
	if err != nil {
		return 0, err
	}
	return n, nil
}

func ParseStringToBool(values url.Values, key string) (bool, error) {
	value := values[key]
	if len(value) == 0 {
		return false, nil
	}
	b, err := strconv.ParseBool(value[0])
	if err != nil {
		return false, err
	}
	return b, nil
}