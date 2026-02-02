package scheduler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(dstart) != 8 {
		return "", errors.New("invalid date format")
	}

	year, err := strconv.Atoi(dstart[:4])
	if err != nil {
		return "", errors.New("invalid date format")
	}
	month, err := strconv.Atoi(dstart[4:6])
	if err != nil || month < 1 || month > 12 {
		return "", errors.New("invalid date format")
	}
	day, err := strconv.Atoi(dstart[6:])
	if err != nil || day < 1 || day > 31 {
		return "", errors.New("invalid date format")
	}

	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	repeat = strings.TrimSpace(repeat)

	if repeat == "y" {
		candidateYear := year

		for {
			candidateYear++
			candidate := time.Date(candidateYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)

			if candidate.Month() != time.Month(month) || candidate.Day() != day {
				candidate = time.Date(candidateYear, time.March, 1, 0, 0, 0, 0, time.UTC)
			}

			if candidate.After(now) {
				return candidate.Format("20060102"), nil
			}
		}
	}

	if len(repeat) >= 2 && repeat[0] == 'd' && repeat[1] == ' ' {
		daysStr := strings.TrimSpace(repeat[2:])
		days, err := strconv.Atoi(daysStr)
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid day count")
		}

		startDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		nextDate := startDate.AddDate(0, 0, days)
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil
	}

	return "", errors.New("unknown repeat rule")
}
