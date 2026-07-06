package task

import (
	"errors"
	"testing"
	"time"
)

func TestTask_Transition_Valid(t *testing.T) {
	task := &Task{ID: "1", Status: StatusPending}

	err := task.Transition(StatusRunning)
	if err != nil {
		t.Errorf("valid transition should not error, got: %v", err)
	}
	if task.Status != StatusRunning {
		t.Errorf("status after transition = %q, want %q", task.Status, StatusRunning)
	}
}

func TestTask_Transition_Invalid(t *testing.T) {
	task := &Task{ID: "1", Status: StatusSucceeded} // 终态

	err := task.Transition(StatusRunning)

	// 1. 必须返回错误
	if err == nil {
		t.Fatal("transition from terminal state should error, got nil")
	}
	// 2. 用 errors.Is 验证错误"身份"确实是 ErrInvalidTransition
	if !errors.Is(err, ErrInvalidTransition) {
		t.Errorf("error is not ErrInvalidTransition: got %v", err)
	}
	// 3. 失败的迁移不能改变原状态
	if task.Status != StatusSucceeded {
		t.Errorf("status changed after invalid transition: %q", task.Status)
	}
}

func TestNewTask_AppliesDefaultStatus(t *testing.T) {
	task, err := NewTask("t1", "echo hi", "")
	if err != nil {
		t.Fatalf("NewTask returned error: %v", err)
	}

	if task.Status != StatusPending {
		t.Errorf("status = %q, want %q", task.Status, StatusPending)
	}
}

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr error
	}{
		{
			name: "valid zero timeout and max retry",
			task: Task{
				ID:      "t1",
				Command: "echo hi",
				Status:  StatusPending,
			},
		},
		{
			name: "valid positive timeout and max retry",
			task: Task{
				ID:       "t1",
				Command:  "echo hi",
				Status:   StatusPending,
				Timeout:  Duration(time.Second),
				MaxRetry: 2,
			},
		},
		{
			name: "missing id",
			task: Task{
				Command: "echo hi",
				Status:  StatusPending,
			},
			wantErr: ErrMissingID,
		},
		{
			name: "missing command",
			task: Task{
				ID:     "t1",
				Status: StatusPending,
			},
			wantErr: ErrMissingCommand,
		},
		{
			name: "negative timeout",
			task: Task{
				ID:      "t1",
				Command: "echo hi",
				Status:  StatusPending,
				Timeout: Duration(-time.Nanosecond),
			},
			wantErr: ErrInvalidTimeout,
		},
		{
			name: "negative max retry",
			task: Task{
				ID:       "t1",
				Command:  "echo hi",
				Status:   StatusPending,
				MaxRetry: -1,
			},
			wantErr: ErrInvalidMaxRetry,
		},
		{
			name: "invalid status",
			task: Task{
				ID:      "t1",
				Command: "echo hi",
				Status:  Status("unknown"),
			},
			wantErr: ErrInvalidTaskStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Validate returned error: %v", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Validate error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
