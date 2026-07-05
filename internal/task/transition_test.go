package task

import "testing"

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		from Status
		to   Status
		want bool
	}{
		// 合法迁移
		{name: "pending to running", from: StatusPending, to: StatusRunning, want: true},
		{name: "pending to canceled", from: StatusPending, to: StatusCanceled, want: true},
		{name: "running to succeeded", from: StatusRunning, to: StatusSucceeded, want: true},
		{name: "running to failed", from: StatusRunning, to: StatusFailed, want: true},
		{name: "running to canceled", from: StatusRunning, to: StatusCanceled, want: true},

		// 非法迁移
		{name: "succeeded cannot go back to running", from: StatusSucceeded, to: StatusRunning, want: false},
		{name: "canceled cannot go to succeeded", from: StatusCanceled, to: StatusSucceeded, want: false},
		{name: "pending cannot skip to succeeded", from: StatusPending, to: StatusSucceeded, want: false},

		// 边界:非法的 from
		{name: "unknown from returns false", from: Status("bogus"), to: StatusRunning, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanTransition(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("CanTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
