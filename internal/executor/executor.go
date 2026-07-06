// Package executor 负责执行任务对应的 shell 命令。
package executor

import (
	"errors"
	"os/exec"
	"time"

	"github.com/mianm12/taskhive/internal/task"
)

// Result 保存一次命令执行的结果。
type Result struct {
	Output   string        // stdout 与 stderr 合并的输出
	ExitCode int           // 进程退出码,0 表示成功
	Duration time.Duration // 执行耗时
	Err      error         // 执行层面的错误(命令没跑起来等);命令跑了但失败不算这里
}

// Run 执行 t.Command 对应的命令,返回执行结果。
func Run(t *task.Task) Result {
	start := time.Now()

	cmd := exec.Command("sh", "-c", t.Command)
	output, err := cmd.CombinedOutput()

	resule := Result{
		Output:   string(output),
		Duration: time.Since(start),
	}

	// 命令成功执行
	if err == nil {
		resule.ExitCode = 0
		return resule
	}

	// 区分两种失败:命令跑了但退出码非 0 vs 命令没跑起来
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		resule.ExitCode = exitErr.ExitCode()
	} else {
		resule.ExitCode = -1
		resule.Err = err
	}

	return resule
}
