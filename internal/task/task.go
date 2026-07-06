// Package task 定义任务领域模型与状态机。
package task

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrMissingID 表示任务 ID 为空。
	ErrMissingID = errors.New("task ID is required")
	// ErrMissingCommand 表示任务命令为空。
	ErrMissingCommand = errors.New("task command is required")
	// ErrInvalidTaskStatus 表示任务状态为非法状态
	ErrInvalidTaskStatus = errors.New("invalid task status")
	// ErrInvalidTimeout 表示任务超时时间为负数。
	ErrInvalidTimeout = errors.New("task timeout must be non-negative")
	// ErrInvalidMaxRetry 表示任务最大重试次数为负数。
	ErrInvalidMaxRetry = errors.New("task max_retry must be non-negative")
)

// Task 是任务领域模型。
type Task struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Command  string        `json:"command"`
	Status   Status        `json:"status"`
	Timeout  time.Duration `json:"timeout"`
	MaxRetry int           `json:"max_retry"`
}

// NewTask 创建一个已应用默认值并通过校验的 Task。
func NewTask(id, command string, status Status) (Task, error) {
	task := Task{
		ID:      id,
		Command: command,
		Status:  status,
	}

	task.ApplyDefaults()

	if err := task.Validate(); err != nil {
		return Task{}, err
	}

	return task, nil
}

// ApplyDefaults 将 Task 中的空值设为默认值。
func (t *Task) ApplyDefaults() {
	if t.Status == "" {
		t.Status = StatusPending
	}
}

// Transition 尝试把任务状态迁移到 to。
// 非法迁移返回包装了 ErrInvalidTransition 的错误。
func (t *Task) Transition(to Status) error {
	if !CanTransition(t.Status, to) {
		return fmt.Errorf("%s → %s: %w", t.Status, to, ErrInvalidTransition)
	}
	t.Status = to
	return nil
}

// Validate 校验 Task 中数据是否符合规范。
func (t *Task) Validate() error {
	if t.ID == "" {
		return ErrMissingID
	}
	if t.Command == "" {
		return ErrMissingCommand
	}
	if t.Timeout < 0 {
		return fmt.Errorf("%w: %s", ErrInvalidTimeout, t.Timeout)
	}
	if t.MaxRetry < 0 {
		return fmt.Errorf("%w: %d", ErrInvalidMaxRetry, t.MaxRetry)
	}
	switch t.Status {
	case StatusCanceled, StatusFailed, StatusPending, StatusRunning, StatusSucceeded:
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrInvalidTaskStatus, t.Status)
	}
}
