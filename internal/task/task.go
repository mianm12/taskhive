// Package task 定义任务领域模型与状态机。
package task

import (
	"fmt"
	"time"
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

// Transition 尝试把任务状态迁移到 to。
// 非法迁移返回包装了 ErrInvalidTransition 的错误。
func (t *Task) Transition(to Status) error {
	if !CanTransition(t.Status, to) {
		return fmt.Errorf("%s → %s: %w", t.Status, to, ErrInvalidTransition)
	}
	t.Status = to
	return nil
}
