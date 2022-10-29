package strings

import (
	"testing"
	"time"
)

func TestParseStringToDate(t *testing.T) {
	m := make(map[string][]string)
	m["start_date"] = []string{"2020-01-06"}

	got, err := ParseStringToDate(m, "start_date")
	if err != nil {
		t.Fatalf("unexpected error=(%v)", err)
	}
	want := time.Date(2020, time.January, 6, 0, 0, 0, 0, time.UTC)

	if *got != want {
		t.Fatalf("ParseStringToDate() want %v, got %v", want, got)
	}
}