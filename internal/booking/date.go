package booking

import (
	"encoding/json"
	"fmt"
	"time"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d Date) ToTime() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

func FromTime(t time.Time) Date {
	return Date{
		Year:  t.Year(),
		Month: int(t.Month()),
		Day:   t.Day(),
	}
}

func NewDate(year, month, day int) Date {
	return Date{Year: year, Month: month, Day: day}
}

func (d Date) Equal(other Date) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	*d = FromTime(parsedTime)
	return nil
}
