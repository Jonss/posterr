package utils

import "time"

// March 25, 2021
const responseLayout = "January 02, 2006"

func ResponseFormatDate(date time.Time) string {
	return date.Format(responseLayout)
}