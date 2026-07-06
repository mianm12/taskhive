package runner

import (
	"fmt"

	"github.com/mianm12/taskhive/internal/executor"
	"github.com/mianm12/taskhive/internal/task"
)

// PrintReport 打印所有任务的执行结果汇总。
func PrintReport(tasks []task.Task, results []executor.Result) {
	var succeeded, failed int

	for i := range tasks {
		r := results[i]
		status := "成功"
		if r.ExitCode != 0 {
			status = "失败"
			failed++
		} else {
			succeeded++
		}

		fmt.Printf("[%s] %-20s 退出码=%-3d 耗时=%-12v %s\n",
			status, tasks[i].Name, r.ExitCode, r.Duration, tasks[i].ID)
	}

	fmt.Printf("\n汇总:共 %d 个,成功 %d,失败 %d\n",
		len(tasks), succeeded, failed)
}
