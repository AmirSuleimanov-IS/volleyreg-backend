package repository

import (
	"math"
	"time"
)

// FullYears — число полных лет на дату at для человека, родившегося birth
func FullYears(birth, at time.Time) int {
	years := at.Year() - birth.Year()
	if at.Month() < birth.Month() || (at.Month() == birth.Month() && at.Day() < birth.Day()) {
		years--
	}
	return years
}

// AverageAge — средний возраст команды на дату соревнования:
// (В1 + В2 + … + Вn) / n, округление до 0,1 года
func AverageAge(birthDates []time.Time, competitionDate time.Time) float64 {
	if len(birthDates) == 0 {
		return 0
	}
	sum := 0
	for _, b := range birthDates {
		sum += FullYears(b, competitionDate)
	}
	avg := float64(sum) / float64(len(birthDates))
	return math.Round(avg*10) / 10
}
