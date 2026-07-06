package executor

import (
	"strings"
	"testing"

	"github.com/mianm12/taskhive/internal/task"
)

func TestRun_Success(t *testing.T) {
	tk := &task.Task{ID: "t1", Command: "echo hello"}
	result := Run(tk)

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
	result := Run(tk)

	if result.ExitCode != 3 {
		t.Errorf("exit code = %d, want 3", result.ExitCode)
	}
	if result.Err != nil {
		t.Errorf("non-zero exit should not set Err, got: %v", result.Err)
	}
}
