package unitime

import (
	"time"

	"github.com/ulikunitz/unixtime"
)

func ToUnixTime(layout, timestr string) (int64, error) {
	t, err := time.Parse(layout, timestr)
	rez := unixtime.Milli(t)
	return int64(float64(rez) / 1000.0), err
}

func timeFromUnix(utime int64) time.Time {
	t := utime * 1000
	return unixtime.FromMilli(t)
}

func FromUnixTime(layout string, utime int64) string {
	t := timeFromUnix(utime)
	return t.Format(layout)
}

func DeltaHours(before, after int64) float64 {
	start, stop := timeFromUnix(before), timeFromUnix(after)
	diff := stop.Sub(start)
	return diff.Hours()
}

func getMidnight(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func UnixTimeToday() int64 {
	midnight := getMidnight(time.Now())
	rez := unixtime.Milli(midnight)
	return int64(float64(rez) / 1000.0)
}

func UnixTimeTomorrow() int64 {
	t := time.Now().Add(time.Hour * 24)
	midnight := getMidnight(t)
	return int64(float64(unixtime.Milli(midnight)) / 1000.0)
}
