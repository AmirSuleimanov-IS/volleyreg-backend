package repository

import (
	"testing"
	"time"
)

func date(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func TestFullYears(t *testing.T) {
	competition := date("2026-11-20")
	cases := []struct {
		birth string
		want  int
	}{
		{"2004-03-14", 22}, // день рождения уже был в этом году
		{"2004-11-21", 21}, // день рождения на следующий день после соревнования
		{"2004-11-20", 22}, // день рождения в день соревнования
		{"2004-12-01", 21}, // день рождения в следующем месяце
	}
	for _, c := range cases {
		if got := FullYears(date(c.birth), competition); got != c.want {
			t.Errorf("FullYears(%s) = %d, want %d", c.birth, got, c.want)
		}
	}
}

func TestAverageAge(t *testing.T) {
	competition := date("2026-11-20")
	births := []time.Time{date("2004-03-14"), date("2004-11-09"), date("2002-05-30")} // 22, 22, 24
	if got := AverageAge(births, competition); got != 22.7 {
		t.Errorf("AverageAge = %v, want 22.7", got)
	}
}

func TestAverageAgeEmpty(t *testing.T) {
	if got := AverageAge(nil, date("2026-11-20")); got != 0 {
		t.Errorf("AverageAge(nil) = %v, want 0", got)
	}
}
