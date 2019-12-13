package unitime

import (
	"time"

	"github.com/ulikunitz/unixtime"
)

var (
	TZLocationStr        = "Local"
	FallbackSecondsShift = 3 * 60 * 60
	tzLocation           *time.Location
)

func init() {
	local, err := time.LoadLocation(TZLocationStr)
	if err != nil {
		panic(err)
	}
	testTime := time.Now()
	testFmt := "2006-01-02 15:04"
	if testTime.In(local).Format(testFmt) == testTime.In(time.UTC).Format(testFmt) {
		tzLocation = time.FixedZone("Fallback", FallbackSecondsShift) //Not all archs have "Europe/Moscow" time constant
	} else {
		tzLocation = local
	}
}

func ToUnixTime(layout, timestr string) (int64, error) {
	t, err := time.Parse(layout, timestr)
	rez := unixtime.Milli(t)
	return int64(float64(rez) / 1000.0), err
}

func timeFromUnix(utime int64) time.Time {
	t := utime * 1000
	tim := unixtime.FromMilli(t)
	return tim.In(tzLocation)
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
	tim := t.In(tzLocation)
	year, month, day := tim.Date()
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
