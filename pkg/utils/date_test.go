package utils

import (
	"testing"
	"time"
)

func TestResponseFormatDate(t *testing.T) {
	date, err := time.Parse("2006-01-02", "2021-03-25")
	if err != nil {
		t.Fatalf("unexpected error parsing date. error=(%v)", err)
	}
	want := "March 25, 2021"
	got := ResponseFormatDate(date)
	if got != want {
		t.Fatalf("ResponseFormatDate() got %v, want %v", got, want)
	}
}
