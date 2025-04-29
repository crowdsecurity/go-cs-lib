package cstime

import (
	"encoding/json"
	"time"
)

// Duration wraps time.Duration and supports parsing with "d" for days.
type Duration time.Duration

// UnmarshalText implements encoding.TextUnmarshaler (used by YAML/JSON libs).
func (d *Duration) UnmarshalText(text []byte) error {
	s := string(text)
	dur, err := ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

// MarshalText returns the standard duration string.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}

// UnmarshalJSON expects a quoted duration string like "2d3h".
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.UnmarshalText([]byte(s))
}

// MarshalJSON outputs a quoted string.
func (d Duration) MarshalJSON() ([]byte, error) {
	s, err := d.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(s))
}

// Set implements the pflag.Value interface.
func (d *Duration) Set(s string) error {
	return d.UnmarshalText([]byte(s))
}

// String returns the standard time.Duration string.
func (d Duration) String() string {
	return time.Duration(d).String()
}

// Type returns the type name for Cobra flags.
func (Duration) Type() string {
	return "duration"
}

// Hours returns the duration as a floating point number of hours.
func (d Duration) Hours() float64 {
	return time.Duration(d).Hours()
}

// Minutes returns the duration as a floating point number of minutes.
func (d Duration) Minutes() float64 {
	return time.Duration(d).Minutes()
}

// Seconds returns the duration as a floating point number of seconds.
func (d Duration) Seconds() float64 {
	return time.Duration(d).Seconds()
}

// Milliseconds returns the duration as an integer millisecond count.
func (d Duration) Milliseconds() int64 {
	return time.Duration(d).Milliseconds()
}

// Nanoseconds returns the duration as an integer nanosecond count.
func (d Duration) Nanoseconds() int64 {
	return time.Duration(d).Nanoseconds()
}
