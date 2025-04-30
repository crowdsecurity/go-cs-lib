package cstime

import (
	"encoding/json"
	"time"
)

// Go's standard library does not support days or weeks as duration units —
// time.ParseDuration understands only hours and smaller (e.g., "72h", "1h30m").
//
// In this package, we define a "day" as exactly 24 hours (24h), and
// optionally support combinations like "2d3h". This simplification assumes
// that days are uniform and does not account for daylight saving time (DST),
// leap seconds, or calendar semantics. This is intentional: we're modeling
// durations, not calendar dates.

// DurationWithDays wraps time.Duration and supports parsing with "d" for days.
type DurationWithDays time.Duration

// UnmarshalText implements encoding.TextUnmarshaler (used by YAML/JSON libs).
func (d *DurationWithDays) UnmarshalText(text []byte) error {
	s := string(text)
	dur, err := ParseDurationWithDays(s)
	if err != nil {
		return err
	}
	*d = DurationWithDays(dur)
	return nil
}

// MarshalText returns the standard duration string.
func (d DurationWithDays) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// UnmarshalJSON expects a quoted duration string like "2d3h".
func (d *DurationWithDays) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.UnmarshalText([]byte(s))
}

// MarshalJSON outputs a quoted string.
func (d DurationWithDays) MarshalJSON() ([]byte, error) {
	s, err := d.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(s))
}

// Set implements the pflag.Value interface.
func (d *DurationWithDays) Set(s string) error {
	return d.UnmarshalText([]byte(s))
}

// String returns the standard time.Duration string.
func (d DurationWithDays) String() string {
	return time.Duration(d).String()
}

// Type returns the type name for Cobra flags.
func (DurationWithDays) Type() string {
	return "duration"
}

// Hours returns the duration as a floating point number of hours.
func (d DurationWithDays) Hours() float64 {
	return time.Duration(d).Hours()
}

// Minutes returns the duration as a floating point number of minutes.
func (d DurationWithDays) Minutes() float64 {
	return time.Duration(d).Minutes()
}

// Seconds returns the duration as a floating point number of seconds.
func (d DurationWithDays) Seconds() float64 {
	return time.Duration(d).Seconds()
}

// Milliseconds returns the duration as an integer millisecond count.
func (d DurationWithDays) Milliseconds() int64 {
	return time.Duration(d).Milliseconds()
}

// Nanoseconds returns the duration as an integer nanosecond count.
func (d DurationWithDays) Nanoseconds() int64 {
	return time.Duration(d).Nanoseconds()
}
