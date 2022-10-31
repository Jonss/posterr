package utils

import (
	"database/sql"
	"time"
)

// March 25, 2021
const responseLayout = "January 2, 2006"

func ResponseFormatDate(date time.Time) string {
	return date.Format(responseLayout)
}

func TimeToNullTime(time time.Time) sql.NullTime {
	return sql.NullTime{Valid: true, Time: time}
}
