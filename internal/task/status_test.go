package task

import "testing"

func TestStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		status Status
		want   bool
	}{
		{name: "pending is not terminal", status: StatusPending, want: false},
		{name: "running is not terminal", status: StatusRunning, want: false},
		{name: "succeeded is terminal", status: StatusSucceeded, want: true},
		{name: "failed is terminal", status: StatusFailed, want: true},
		{name: "canceled is terminal", status: StatusCanceled, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsTerminal()
			if got != tt.want {
				t.Errorf("Status(%q).IsTerminal() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
