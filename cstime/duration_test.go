package cstime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/crowdsecurity/go-cs-lib/cstest"
)

func TestUnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr string
	}{
		{
			name:  "HoursMinutes",
			input: "1h30m",
			want:  time.Hour + 30*time.Minute,
		}, {
			name:  "Days",
			input: "2d",
			want:  48 * time.Hour,
		}, {
			name:    "TooManyDays",
			input:   "2d4d",
			wantErr: `invalid duration "2d4d": 'd' can appear only as the first unit`,
		}, {
			name:    "Weeks",
			input:   "1w",
			want:    7 * 24 * time.Hour,
			wantErr: `time: unknown unit "w" in duration "1w"`,
		}, {
			name:  "MixedUnits",
			input: "2d3h4m5s",
			want:  2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second,
		}, {
			name:    "FractionalDay",
			input:   "1.5d",
			wantErr: `invalid day value in duration "1.5d"`,
		}, {
			name:    "Invalid",
			input:   "invalid",
			wantErr: "invalid",
		}, {
			name:    "WrongUnitOrder",
			input:   "1h30m2d",
			wantErr: `invalid day value in duration "1h30m2d"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var d DurationWithDays
			err := d.UnmarshalText([]byte(tc.input))
			cstest.RequireErrorContains(t, err, tc.wantErr)
			if tc.wantErr != "" {
				return
			}
			assert.Equal(t, tc.want, time.Duration(d))
		})
	}
}

func TestMarshalText(t *testing.T) {
	tests := []struct {
		name string
		d    DurationWithDays
		want string
	}{
		{"Zero", DurationWithDays(0), "0s"},
		{"Simple", DurationWithDays(2*time.Hour + 30*time.Minute), "2h30m0s"},
		{"OneDay", DurationWithDays(24 * time.Hour), "24h0m0s"},
		{"OneWeek", DurationWithDays(7 * 24 * time.Hour), "168h0m0s"},
		{"Complex", DurationWithDays((7*24+2)*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second), "173h4m5s"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.d.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(b))
		})
	}
}

func TestDurationFlagParsing(t *testing.T) {
	var dur DurationWithDays

	cmd := &cobra.Command{}
	cmd.Flags().Var(&dur, "duration", "custom duration with day support")

	err := cmd.ParseFlags([]string{"--duration=2d3h"})
	require.NoError(t, err)

	want := 2*24*time.Hour + 3*time.Hour
	assert.Equal(t, DurationWithDays(want), dur)

	err = cmd.ParseFlags([]string{"--duration=3h2d"})
	cstest.RequireErrorContains(t, err, `invalid argument "3h2d" for "--duration" flag: invalid day value in duration "3h2d"`)
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	type Config struct {
		Timeout DurationWithDays `json:"timeout"`
	}

	jsonInput := `{"timeout": "2d3h45m"}`

	var cfg Config
	err := json.Unmarshal([]byte(jsonInput), &cfg)
	require.NoError(t, err)

	want := 2*24*time.Hour + 3*time.Hour + 45*time.Minute
	assert.Equal(t, DurationWithDays(want), cfg.Timeout)
}

func TestDuration_UnmarshalYAML(t *testing.T) {
	type Config struct {
		Timeout DurationWithDays `yaml:"timeout"`
	}

	yamlInput := "timeout: 2d3h45m"

	var cfg Config
	err := yaml.Unmarshal([]byte(yamlInput), &cfg)
	require.NoError(t, err)

	want := 2*24*time.Hour + 3*time.Hour + 45*time.Minute
	assert.Equal(t, DurationWithDays(want), cfg.Timeout)
}

func TestDuration_Methods(t *testing.T) {
	d := DurationWithDays(2*time.Hour + 30*time.Minute)

	assert.Equal(t, 2.5, d.Hours())      //nolint:testifylint
	assert.Equal(t, 150.0, d.Minutes())  //nolint:testifylint
	assert.Equal(t, 9000.0, d.Seconds()) //nolint:testifylint
	assert.Equal(t, int64(9000000), d.Milliseconds())
	assert.Equal(t, int64(9000000000000), d.Nanoseconds())
}
