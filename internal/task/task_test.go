package task

import (
	"errors"
	"testing"
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
