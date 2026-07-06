package executor

import (
	"strings"
	"testing"
	"time"

	"github.com/mianm12/taskhive/internal/task"
)

func TestRun_Success(t *testing.T) {
	tk := &task.Task{ID: "t1", Command: "echo hello"}
	result := Run(tk, Config{})

	if result.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", result.ExitCode)
	}
	if !strings.Contains(result.Output, "hello") {
		t.Errorf("output = %q, want to contain 'hello'", result.Output)
	}
	if result.Err != nil {
		t.Errorf("unexpected err: %v", result.Err)
	}
}

func TestRun_NonZeroExitCode(t *testing.T) {
	tk := &task.Task{ID: "t2", Command: "exit 3"}
	result := Run(tk, Config{})

	if result.ExitCode != 3 {
		t.Errorf("exit code = %d, want 3", result.ExitCode)
	}
	if result.Err != nil {
		t.Errorf("non-zero exit should not set Err, got: %v", result.Err)
	}
}

func TestRun_Timeout(t *testing.T) {
	t.Skip("超时杀进程在 Linux 上需杀整个进程组,依赖阶段 3 并发方案,见 docs/notes/TODO.md")

	tk := &task.Task{ID: "t3", Command: "sleep 5"}
	cfg := Config{Timeout: 200 * time.Millisecond}

	result := Run(tk, cfg)

	if !result.TimedOut {
		t.Errorf("expected TimedOut=true, got false")
	}
	if result.Duration >= 5*time.Second {
		t.Errorf("command was not killed by timeout, took %v", result.Duration)
	}
}

func TestRun_RetryThenGiveUp(t *testing.T) {
	tk := &task.Task{ID: "t4", Command: "exit 1"}
	cfg := Config{MaxRetries: 2, RetryDelay: 10 * time.Millisecond}

	result := Run(tk, cfg)

	if result.ExitCode != 1 {
		t.Errorf("exit code = %d, want 1", result.ExitCode)
	}
	if result.Attempts != 3 { // 1 次初始 + 2 次重试
		t.Errorf("attempts = %d, want 3", result.Attempts)
	}
}
