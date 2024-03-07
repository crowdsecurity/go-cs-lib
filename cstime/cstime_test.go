package cstime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     time.Duration
		wantErr  bool
	}{
		{
			name:     "only days",
			input:    "2d",
			want:     48 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "only days, negative",
			input:    "-2d",
			want:     -48 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "days and hours",
			input:    "1d2h",
			want:     26 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "days and hours, negative",
			input:    "-1d2h",
			want:     -26 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "hours minutes seconds",
			input:    "2h45m30s",
			want:     2*time.Hour + 45*time.Minute + 30*time.Second,
			wantErr:  false,
		},
		{
			name:     "invalid format",
			input:    "2x",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "invalid days format",
			input:    "2x3h",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "invalid remainder format",
			input:    "2d3x",
			want:     0,
			wantErr:  true,
		},
		{
			name:     "empty input",
			input:    "",
			want:     0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
