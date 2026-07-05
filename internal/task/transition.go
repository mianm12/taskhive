// Package task 定义任务领域模型与状态机。
package task

import "errors"

// ErrInvalidTransition 表示尝试了一次非法的状态迁移。
var ErrInvalidTransition = errors.New("invalid status transition")

// validTransitions 定义每个状态允许迁移到哪些状态。
// key 是当前状态,value 是允许的目标状态集合。
var validTransitions = map[Status][]Status{
	StatusPending:   {StatusRunning, StatusCanceled},                 // 等待 → 运行 / 取消
	StatusRunning:   {StatusSucceeded, StatusFailed, StatusCanceled}, // 运行 → 成功 / 失败 / 取消
	StatusSucceeded: {},                                              // 终态
	StatusFailed:    {},                                              // 终态
	StatusCanceled:  {},                                              // 终态
}

// CanTransition 报告从 from 迁移到 to 是否合法。
func CanTransition(from, to Status) bool {
	// 对于非法的 from 状态，vaidTransitions[from] 将返回 nil，for 循环不会执行，最终返回 false。
	for _, allowed := range validTransitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}
