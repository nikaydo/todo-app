package date

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrDaysOutOfRange     = errors.New("specified days exceed the number of days in a month")
	ErrWeekDaysOutOfRange = errors.New("specified day of the week or month is invalid")
	ErrTooManyDays        = errors.New("too many days specified")
	ErrInvalidRepeat      = errors.New("invalid repeat rule specified")
	ErrNegativeDay        = errors.New("negative day value is too large")
)

const (
	layout = "20060102"
)

func parseIntList(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	var result []int
	for _, part := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}

func month(now time.Time, startStr string, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) < 2 {
		return "", ErrInvalidRepeat
	}
	dayList, err := parseIntList(parts[1])
	if err != nil {
		return "", err
	}
	start, err := time.Parse(layout, startStr)
	if err != nil {
		return "", err
	}
	monthFilter, err := filterMonth(parts)
	if err != nil {
		return "", err
	}
	date := start
	if !date.After(now) {
		date = now.AddDate(0, 0, 1)
	}
	for {
		daysInCurrentMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()

		isValidDay, err := isValidDay(dayList, date, daysInCurrentMonth)
		if err != nil {
			return "", err
		}
		isValidMonth := len(monthFilter) == 0 || monthFilter[int(date.Month())]

		if isValidDay && isValidMonth {
			return date.Format(layout), nil
		}
		date = date.AddDate(0, 0, 1)
	}
}

func filterMonth(parts []string) (map[int]bool, error) {
	monthFilter := make(map[int]bool)
	var monthList []int
	var err error
	if len(parts) == 3 {
		monthList, err = parseIntList(parts[2])
		if err != nil {
			return monthFilter, err
		}
	}
	for _, m := range monthList {
		if m < 1 || m > 12 {
			return monthFilter, ErrWeekDaysOutOfRange
		}
		monthFilter[m] = true
	}
	return monthFilter, nil
}

func isValidDay(dayList []int, date time.Time, daysInCurrentMonth int) (bool, error) {
	for _, d := range dayList {
		switch {
		case d > 0 && d <= 31 && date.Day() == d:
			return true, nil
		case d == -1 && date.Day() == daysInCurrentMonth:
			return true, nil
		case d == -2 && date.Day() == daysInCurrentMonth-1:
			return true, nil
		case d > 31 || d < -2:
			return true, ErrWeekDaysOutOfRange
		}
	}
	return false, nil
}

func week(now time.Time, startStr string, repeat string) (string, error) {
	parts := strings.Split(repeat, " ")
	if len(parts) != 2 {
		return "", ErrInvalidRepeat
	}
	days, err := parseIntList(parts[1])
	if err != nil {
		return "", err
	}
	start, err := time.Parse(layout, startStr)
	if err != nil {
		return "", err
	}
	weekdays := make(map[int]bool)
	for _, d := range days {
		if d < 1 || d > 7 {
			return "", ErrWeekDaysOutOfRange
		}
		weekdays[d] = true
	}
	date := start
	if !date.After(now) {
		date = now.AddDate(0, 0, 1)
	}
	for {
		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7
		}
		if weekdays[wd] {
			return date.Format(layout), nil
		}
		date = date.AddDate(0, 0, 1)
	}
}

func day(now time.Time, dstart string, repeat string) (string, error) {
	r := strings.Split(repeat, " ")
	if len(r) != 2 {
		return "", ErrInvalidRepeat
	}
	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", err
	}
	n, err := strconv.Atoi(r[1])
	if err != nil {
		return "", ErrInvalidRepeat
	}
	if n >= 400 {
		return "", ErrTooManyDays
	}
	if n == 1 {
		return date.AddDate(0, 0, 1).Format(layout), nil
	}
	for {
		date = date.AddDate(0, 0, n)
		if date.After(now) {
			return date.Format(layout), err
		}
	}
}

func year(now time.Time, dstart string) (string, error) {
	start, err := time.Parse(layout, dstart)
	if err != nil {
		return "", err
	}
	for {
		start = start.AddDate(1, 0, 0)
		if start.After(now) {
			return start.Format(layout), nil
		}
	}
}

func NextDate(now time.Time, dateStart string, repeat string) (string, error) {
	r := strings.Split(repeat, " ")

	switch r[0] {
	case "d":
		return day(now, dateStart, repeat)
	case "y":
		return year(now, dateStart)
	case "w":
		return week(now, dateStart, repeat)
	case "m":
		return month(now, dateStart, repeat)
	default:
		return "", ErrInvalidRepeat
	}
}
