// Package command 定义 taskhive 的 CLI 命令树。
package command

import "github.com/spf13/cobra"

// NewRootCmd 构造根命令并挂载所有子命令。
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "taskhive",
		Short: "TaskHive 是一个任务调度与执行工具",
	}

	root.AddCommand(newRunCmd())
	root.AddCommand(newVersionCmd())
	return root
}
