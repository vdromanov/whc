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

func ToUnixTime(t time.Time) int64 {
	ut := unixtime.Milli(t)
	return int64(float64(ut) / 1000.0)
}

func FormatToUnixTime(layout, timestr string) (int64, error) {
	t, err := time.Parse(layout, timestr)
	return ToUnixTime(t), err

}

func timeFromUnix(utime int64) time.Time {
	t := utime * 1000
	tim := unixtime.FromMilli(t)
	return tim.In(tzLocation)
}

func FormatFromUnixTime(layout string, utime int64) string {
	t := timeFromUnix(utime)
	return t.Format(layout)
}

func DeltaHoursUnixTime(before, after int64) float64 {
	start, stop := timeFromUnix(before), timeFromUnix(after)
	diff := stop.Sub(start)
	return diff.Hours()
}

func GetBeginningOfDay(fmtStr, valuesStr string) (time.Time, error) {
	t, err := time.Parse(fmtStr, valuesStr)
	tim := t.In(tzLocation)
	year, month, day := tim.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, tzLocation), err
}

func NextDay(t time.Time) time.Time {
	return t.In(tzLocation).Add(time.Hour * 24)
}
