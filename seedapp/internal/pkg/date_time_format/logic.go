package date_time_format

import (
	"errors"
	"seedapp/internal/adapter/repository/sqlx/querier"
	"strconv"
	"strings"
	"time"
)

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func GetIntDate(date string) (int, int, int, error) {
	yr, err := strconv.Atoi(date[0:4])
	if err != nil {
		return 0, 0, 0, err
	}

	mth, err := strconv.Atoi(date[4:6])
	if err != nil {
		return 0, 0, 0, err
	}

	dt, err := strconv.Atoi(date[len(date)-2:])
	if err != nil {
		return 0, 0, 0, err
	}

	return yr, mth, dt, nil
}

func RangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

func ParseDate(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(DATE_NUMERIC, dateStr)
	if err != nil {
		return time.Now(), err
	}
	return parsedDate, nil
}

func ParseDateTime(dateStr string) (time.Time, error) {
	parsedDate, err := time.Parse(DATETIME_NUMERIC, dateStr)
	if err != nil {
		return time.Now(), err
	}
	return parsedDate, nil
}

func ParseDate_ISO8601(dateStr string) (time.Time, error) {
	if !strings.Contains(dateStr, "-") {
		dateStr = dateStr[:4] + "-" + dateStr[4:6] + "-" + dateStr[6:]
	}
	parsedDate, err := time.Parse(DATE_ISO_8601, dateStr)
	if err != nil {
		return time.Now(), err
	}
	return parsedDate, nil
}

func ParseDateTime_ISO8601(dateStr string) (time.Time, error) {
	if !strings.Contains(dateStr, "-") && !strings.Contains(dateStr, ":") {
		dateStr = dateStr[:4] + "-" + dateStr[4:6] + "-" + dateStr[6:8] + " " + dateStr[8:10] + ":" + dateStr[10:12] + ":" + dateStr[12:]
	}

	parsedDate, err := time.Parse(DATETIME_ISO_8601, dateStr)
	if err != nil {
		return time.Now(), err
	}
	return parsedDate, nil
}

func ParseAndFormatDateTime(date time.Time) (string, error) {
	parsedDate, err := time.Parse(DATE_LAYOUT, date.String())
	if err != nil {
		return "", err
	}

	return parsedDate.Format(querier.DATE_FORMAT), nil
}

func ParseDateToString(date time.Time) string {
	if !date.IsZero() {
		return date.Format(DATE_NUMERIC)
	}
	return ""
}

// parse to comply DB date format
// parse time.RFC3339 to DB layout. See : querier.DATE_FORMAT
func ParseAndFormatDate(dateTime time.Time) (string, error) {
	if !dateTime.IsZero() {
		parsedDate, err := time.Parse(DATE_LAYOUT, dateTime.String())
		if err != nil {
			return "", err
		}

		return parsedDate.Format(querier.DATE_FORMAT), nil
	}

	return "", nil
}
func ParseStringDateUTCToStringDateNumeric(dateStr string) (string, error) {
	t, err := time.Parse(DATE_LAYOUT, dateStr)
	if err != nil {
		return "", err
	}

	return t.Format(DATE_NUMERIC), nil
}
func ParseStringDateISO8601ToStringDateNumeric(dateStr string) (string, error) {
	t, err := time.Parse(DATE_ISO_8601, dateStr[:len(DATE_ISO_8601)])
	if err != nil {
		return "", err
	}

	return t.Format(DATE_NUMERIC), nil
}

func ParseStringLocalTimeZoneToStringUTC(dateStr, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	t, err := time.ParseInLocation(DATE_ISO_8601, dateStr, loc)
	if err != nil {
		return "", err
	}

	return t.UTC().Format(DATE_LAYOUT), nil
}

func ParseUTCToStringLocalTimeZone(dateStr string, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	t, err := time.Parse(DATE_LAYOUT, dateStr)
	if err != nil {
		return "", err
	}

	return t.In(loc).Format(DATE_ISO_8601), nil
}

func DaysBetweenValidated(start, end time.Time) (int, error) {
	if start.After(end) {
		return 0, errors.New("start date cannot be after end date")
	}
	return DaysBetween(start, end), nil
}

func DaysBetween(start, end time.Time) int {
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())
	dur := e.Sub(s)
	return int(dur.Hours() / 24)
}
