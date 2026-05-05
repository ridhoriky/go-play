package validation

import (
	"fmt"
	"time"
)

type DateRange struct {
	From time.Time
	To   time.Time
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func Parse(period, dateFrom, dateTo string) (DateRange, error) {
	now := time.Now()

	switch period {
	case "today", "":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.Add(24*time.Hour - time.Second)
		return DateRange{From: start, To: end}, nil

	case "this_week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday = 7
		}
		start := now.AddDate(0, 0, -(weekday - 1))
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 0, 6)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, now.Location())
		return DateRange{From: start, To: end}, nil

	case "this_month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, now.Location())
		return DateRange{From: start, To: end}, nil

	case "custom":
		return parseCustomRange(dateFrom, dateTo)

	default:
		return DateRange{}, fmt.Errorf("invalid period: %s", period)
	}
}

func parseCustomRange(dateFrom, dateTo string) (DateRange, error) {
	if dateFrom == "" || dateTo == "" {
		return DateRange{}, &ValidationError{
			Message: "date_from and date_to are required when period is 'custom'",
		}
	}

	const layout = "2006-01-02"

	from, err := time.Parse(layout, dateFrom)
	if err != nil {
		return DateRange{}, &ValidationError{
			Message: "date_from format must be YYYY-MM-DD",
		}
	}

	to, err := time.Parse(layout, dateTo)
	if err != nil {
		return DateRange{}, &ValidationError{
			Message: "date_to format must be YYYY-MM-DD",
		}
	}

	if to.Before(from) {
		return DateRange{}, &ValidationError{
			Message: "date_to must be after date_from",
		}
	}

	// set to end of day
	to = time.Date(to.Year(), to.Month(), to.Day(), 23, 59, 59, 0, to.Location())

	return DateRange{From: from, To: to}, nil
}
