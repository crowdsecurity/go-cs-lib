package cstime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/crowdsecurity/go-cs-lib/cstest"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr string
	}{
		{
			name:  "OnlyDays",
			input: "2d",
			want:  48 * time.Hour,
		},
		{
			name:  "NegativeDays",
			input: "-2d",
			want:  -48 * time.Hour,
		},
		{
			name:  "DaysHours",
			input: "1d2h",
			want:  26 * time.Hour,
		},
		{
			name:  "NegativeDaysHours",
			input: "-1d2h",
			want:  -26 * time.Hour,
		},
		{
			name:    "MisplacedNegative",
			input:   "1d-2h",
			wantErr: "misplaced negative sign",
		},
		{
			name:  "HoursMinutesSeconds",
			input: "2h45m30s",
			want:  2*time.Hour + 45*time.Minute + 30*time.Second,
		},
		{
			name:    "InvalidUnit",
			input:   "2x",
			wantErr: `time: unknown unit "x" in duration "2x"`,
		},
		{
			name:    "InvalidUnitDays",
			input:   "2x3h",
			wantErr: `time: unknown unit "x" in duration "2x3h"`,
		},
		{
			name:    "InvalidUnitRemainder",
			input:   "2d3x",
			wantErr: `time: unknown unit "x" in duration "3x"`,
		},
		{
			name:    "EmptyInput",
			input:   "",
			wantErr: "empty duration string",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseDurationWithDays(tc.input)
			cstest.RequireErrorContains(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
