// Package version 提供构建期注入的版本信息。
package version

import "fmt"

// 这三个变量在构建时通过 -ldflags "-X" 注入,
// 不注入时保持下面的默认值。
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

// String 返回一行人类可读的版本摘要。
func String() string {
	return format(Version, Commit, BuildDate)
}

// format 是纯函数,方便写测试——
// 直接测 String() 要去改包级全局变量,不干净。
func format(version, commit, buildDate string) string {
	return fmt.Sprintf("taskhive %s (commit %s, built %s)", version, commit, buildDate)
}
