// Package executor 负责执行任务对应的 shell 命令。
package executor

import (
	"context"
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
	Attempts int           // 实际执行次数(含重试)
	TimedOut bool          // 是否因超时被终止
	Err      error         // 执行层错误(命令没跑起来等)
}

// Config 控制一次执行的行为。
type Config struct {
	Timeout    time.Duration // 单次执行超时;0 表示不限时
	MaxRetries int           // 失败后最多重试几次;0 表示不重试
	RetryDelay time.Duration // 每次重试前的等待
}

// Run 执行 t.Command,按 cfg 应用超时与重试,返回最终结果。
func Run(t *task.Task, cfg Config) Result {
	var result Result

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		result = runOnce(t, cfg.Timeout)
		result.Attempts = attempt + 1

		if result.ExitCode == 0 {
			return result // 成功执行,直接返回
		}

		// 还有重试机会就等一下再来
		if attempt < cfg.MaxRetries {
			time.Sleep(cfg.RetryDelay)
		}
	}

	return result
}

// runOnce 执行单次命令,应用超时控制。
func runOnce(t *task.Task, timeout time.Duration) Result {
	// TODO(stage-3): sh -c 启动的子进程在 Linux 上不随父进程被杀，
	// 超时对 sleep 这类会失效。需改用进程组(Setpgid) + syscall.Kill(-pid)
	// + goroutine/select 控制,依赖阶段 3 并发知识。
	start := time.Now()

	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", t.Command)
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

	// 先判断是不是超时导致的
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		resule.TimedOut = true
		resule.ExitCode = -1
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
