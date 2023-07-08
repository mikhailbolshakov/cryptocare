package kit

import (
	"errors"
	"time"
)

func MillisFromTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func TimeFromMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

var timeLayout = "15:04"

// Time - time (hour:min) representation in format 15:04
type Time struct {
	time.Time
}

// TimeRange represents time interval in format [15:00, 18:00]
type TimeRange [2]Time

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(timeLayout) + `"`), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	// len(`"23:59"`) == 7
	if len(s) != 7 {
		return errors.New("time parsing error")
	}
	ret, err := time.Parse(timeLayout, s[1:6])
	if err != nil {
		return err
	}
	t.Time = ret
	return nil
}

func (t *Time) Parse(s string) error {
	tm, err := time.Parse(timeLayout, s)
	if err != nil {
		return err
	}
	t.Time = tm
	return nil
}

func (t *Time) String() string {
	return t.Format(timeLayout)
}

func (t TimeRange) MustParse(s1, s2 string) TimeRange {
	tFrom, err := time.Parse(timeLayout, s1)
	if err != nil {
		panic(err)
	}
	tTo, err := time.Parse(timeLayout, s2)
	if err != nil {
		panic(err)
	}
	return TimeRange{Time{Time: tFrom}, Time{Time: tTo}}
}

func (t TimeRange) Valid() bool {
	return t[0].Before(t[1].Time)
}

// Now is the current time
func Now() time.Time {
	return time.Now().Round(time.Microsecond).UTC()
}

// NowNano is the current time in UNIX NANO format
func NowNanos() int64 {
	return time.Now().UTC().UnixNano()
}

// NowMillis is the current time in millis
func NowMillis() int64 {
	return Millis(time.Now().UTC())
}

// Millis is a convenience method to get milliseconds since epoch for provided Time.
func Millis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// Diff properly calculates difference between two dates in year, month etc.
func Diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}

// Date returns a date of a passed timestamp without time
func Date(date time.Time) time.Time {
	y, m, d := date.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// NowDate returns current date without time
func NowDate() time.Time {
	return Date(Now())
}
