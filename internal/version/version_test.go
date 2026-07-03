package version

import "testing"

func Test_format(t *testing.T) {
	tests := []struct {
		name string // 此测试用例的描述
		// 目标函数的命名输入参数。
		version   string
		commit    string
		buildDate string
		want      string
	}{
		{
			name:      "默认值",
			version:   "dev",
			commit:    "none",
			buildDate: "unknown",
			want:      "taskhive dev (commit none, built unknown)",
		},
		{
			name:      "注入后的自定义值",
			version:   "v0.0.1",
			commit:    "abc1234",
			buildDate: "2026-07-03T00:00:00Z",
			want:      "taskhive v0.0.1 (commit abc1234, built 2026-07-03T00:00:00Z)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := format(tt.version, tt.commit, tt.buildDate)
			if got != tt.want {
				t.Errorf("format() = %q, want %q", got, tt.want)
			}
		})
	}
}
