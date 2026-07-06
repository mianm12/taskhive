package command

import (
	"fmt"

	"github.com/mianm12/taskhive/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "打印版本信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.String())
		},
	}
}
