package command

import (
	"fmt"

	"github.com/mianm12/taskhive/internal/runner"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "从 JSON 文件加载并串行执行任务",
		RunE: func(cmd *cobra.Command, args []string) error {
			tasks, err := runner.LoadTasks(file)
			if err != nil {
				return err
			}

			fmt.Printf("加载了 %d 个任务,开始执行...\n\n", len(tasks))
			results := runner.RunAll(tasks)
			runner.PrintReport(tasks, results)
			return nil
		},
	}

	// -f / --file 标志,必填
	cmd.Flags().StringVarP(&file, "file", "f", "", "任务 JSON 文件路径(必填)")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
