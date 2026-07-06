package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mianm12/taskhive/internal/task"
)

func TestLoadTasks(t *testing.T) {
	// 在测试专用的临时目录里造一个 JSON 文件
	dir := t.TempDir()
	path := filepath.Join(dir, "tasks.json")
	content := `[
		{"id": "t1", "name": "echo", "command": "echo hi"},
		{"id": "t2", "name": "fail", "command": "exit 1", "max_retry": 2}
	]`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("写测试文件失败: %v", err)
	}

	tasks, err := LoadTasks(path)
	if err != nil {
		t.Fatalf("LoadTasks 出错: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("任务数 = %d, want 2", len(tasks))
	}
	if tasks[0].ID != "t1" || tasks[0].Command != "echo hi" {
		t.Errorf("task[0] 解析错误: %+v", tasks[0])
	}
	if tasks[1].MaxRetry != 2 {
		t.Errorf("task[1].MaxRetry = %d, want 2", tasks[1].MaxRetry)
	}
}

func TestLoadTasks_FileNotFound(t *testing.T) {
	_, err := LoadTasks("/no/such/file.json")
	if err == nil {
		t.Fatal("读取不存在的文件应报错,却得到 nil")
	}
}

func TestLoadTasks_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{ this is not valid json"), 0o644); err != nil {
		t.Fatalf("写测试文件失败: %v", err)
	}

	_, err := LoadTasks(path)
	if err == nil {
		t.Fatal("解析非法 JSON 应报错,却得到 nil")
	}
}

func TestRunAll(t *testing.T) {
	tasks := []task.Task{
		{ID: "t1", Name: "成功", Command: "echo ok"},
		{ID: "t2", Name: "失败", Command: "exit 1"},
	}

	results := RunAll(tasks)

	if len(results) != 2 {
		t.Fatalf("结果数 = %d, want 2", len(results))
	}
	if results[0].ExitCode != 0 {
		t.Errorf("task1 退出码 = %d, want 0", results[0].ExitCode)
	}
	if results[1].ExitCode != 1 {
		t.Errorf("task2 退出码 = %d, want 1", results[1].ExitCode)
	}
}
