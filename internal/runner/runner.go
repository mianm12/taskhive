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
		if err := tasks[i].Transition(task.StatusRunning); err != nil {
			// FIXME(stage-2): 迁移失败时命令并未执行，用 executor.Result 表达"未执行"
			// 是语义错配，Err 也未被报表消费。待阶段 2 有完整任务模型后重构结果表达。
			res := executor.Result{
				ExitCode: -1,
				Err:      err,
			}
			results = append(results, res)

			continue
		}

		cfg := executor.Config{
			Timeout:    tasks[i].Timeout.Std(),
			MaxRetries: tasks[i].MaxRetry,
			RetryDelay: 500 * time.Millisecond,
		}
		res := executor.Run(&tasks[i], cfg)

		nextStatus := task.StatusFailed
		if res.ExitCode == 0 {
			nextStatus = task.StatusSucceeded
		}

		// 此处 Status 必为 Running(上面迁移失败已 continue),
		// Running → 终态在状态机中一定合法。若这里失败,说明不变量被破坏,
		// 是代码 bug 而非可预期错误,直接 panic 以便立即暴露。
		if err := tasks[i].Transition(nextStatus); err != nil {
			panic(fmt.Sprintf("unreachable: running task cannot transition to %s: %v", nextStatus, err))
		}

		results = append(results, res)
	}
	return results
}
