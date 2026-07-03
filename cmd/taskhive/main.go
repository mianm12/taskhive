package main

import (
	"fmt"
	"os"

	"github.com/mianm12/taskhive/internal/version"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version.String())
		return
	}
	fmt.Println("用法: taskhive version")
}
