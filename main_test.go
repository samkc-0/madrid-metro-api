package main

import (
	"testing"
)

func TestLoadStops(t *testing.T) {
	stops, err := loadStops("data/stops.txt")

	if err != nil {
		t.Fatal(err)
	}

	if len(stops) == 0 {
		t.Fatal("no stops found")
	}

	firstStop := stops[0]
	want := "PLAZA DE CASTILLA"
	got := firstStop.Name
	if got != "PLAZA DE CASTILLA" {
		t.Fatalf("expected first stop to be %s, got %s", want, got)
	}
}
