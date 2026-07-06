package task

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    time.Duration
		wantErr error
	}{
		{
			name: "valid duration string",
			data: `"1h30m"`,
			want: time.Hour + 30*time.Minute,
		},
		{
			name:    "duration must be string",
			data:    `30`,
			wantErr: ErrDurationNotString,
		},
		{
			name:    "invalid duration string",
			data:    `"soon"`,
			wantErr: ErrInvalidDuration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Duration

			err := json.Unmarshal([]byte(tt.data), &got)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UnmarshalJSON error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("UnmarshalJSON returned error: %v", err)
			}
			if got.Std() != tt.want {
				t.Errorf("Duration.Std() = %s, want %s", got.Std(), tt.want)
			}
		})
	}
}
