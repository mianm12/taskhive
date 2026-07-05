// Package task 定义任务领域模型与状态机。
package task

// Status 表示任务的生命周期状态。
type Status string

// 任务生命周期的所有合法状态。
const (
	StatusPending   Status = "pending"   // 等待执行
	StatusRunning   Status = "running"   // 执行中
	StatusSucceeded Status = "succeeded" // 执行成功
	StatusFailed    Status = "failed"    // 执行失败
	StatusCanceled  Status = "canceled"  // 已取消
)

// IsTerminal 报告该状态是否为终态(不可再迁移到其他状态)。
func (s Status) IsTerminal() bool {
	switch s {
	case StatusSucceeded, StatusFailed, StatusCanceled:
		return true
	default:
		return false
	}
}
