package main

import (
	"testing"
	"time"
)

func TestProcessDateCommand(t *testing.T) {
	tests := []struct {
		Command      string
		ExpectedFrom time.Time
		ExpectedTo   time.Time
	}{
		{
			Command:      "3/3/1986",
			ExpectedFrom: time.Date(1986, 3, 3, 0, 0, 0, 0, time.UTC),
			ExpectedTo:   time.Date(1986, 3, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			Command:      "April 11, 2016",
			ExpectedFrom: time.Date(2016, 4, 11, 0, 0, 0, 0, time.UTC),
			ExpectedTo:   time.Date(2016, 4, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			Command:      "8/1",
			ExpectedFrom: time.Date(time.Now().Year(), 8, 1, 0, 0, 0, 0, time.UTC),
			ExpectedTo:   time.Date(time.Now().Year(), 8, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Command:      "August 1",
			ExpectedFrom: time.Date(time.Now().Year(), 8, 1, 0, 0, 0, 0, time.UTC),
			ExpectedTo:   time.Date(time.Now().Year(), 8, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Command:      "2007",
			ExpectedFrom: time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedTo:   time.Date(2008, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Command:      "today",
			ExpectedFrom: time.Now().Truncate(24 * time.Hour),
			ExpectedTo:   time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1),
		},
		{
			Command:      "yesterday",
			ExpectedFrom: time.Now().Truncate(24*time.Hour).AddDate(0, 0, -1),
			ExpectedTo:   time.Now().Truncate(24 * time.Hour),
		},
	}

	for _, test := range tests {
		from, to, err := ProcessDate(test.Command)
		if !from.Equal(test.ExpectedFrom) {
			t.Fatalf("command: %s expected from to be %v, got %v", test.Command, test.ExpectedFrom, from)
		}
		if !to.Equal(test.ExpectedTo) {
			t.Fatalf("command: %s expected from to be %v, got %v", test.Command, test.ExpectedTo, to)
		}
		if err != nil {
			t.Fatalf("command: %s expected no error, got %v", test.Command, err)
		}
	}
}
