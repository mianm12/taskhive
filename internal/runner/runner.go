// Package runner 编排从文件加载任务并串行执行的流程。
package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mianm12/taskhive/internal/executor"
	"github.com/mianm12/taskhive/internal/task"
)

// LoadTasks 从 JSON 文件读取任务列表。
func LoadTasks(path string) ([]task.Task, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading task file %s: %w", path, err)
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("parse the task JSON: %w", err)
	}

	for i := range tasks {
		// 反序列化的 Task 数据没有经过 NewTask() 初始化，这里手动设置默认值
		tasks[i].ApplyDefaults()

		if err := tasks[i].Validate(); err != nil {
			return nil, fmt.Errorf("task[%d] (id=%q) invalid: %w", i, tasks[i].ID, err)
		}
	}

	return tasks, nil
}

// RunAll 串行执行所有任务,返回每个任务的执行结果。
func RunAll(tasks []task.Task) []executor.Result {
	results := make([]executor.Result, 0, len(tasks))
	for i := range tasks {
		cfg := executor.Config{
			Timeout:    tasks[i].Timeout,
			MaxRetries: tasks[i].MaxRetry,
			RetryDelay: 500 * time.Millisecond,
		}
		res := executor.Run(&tasks[i], cfg)
		results = append(results, res)
	}
	return results
}
