// Package task 定义任务领域模型与状态机。
package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrDurationNotString 表示 JSON 中的 Duration 不是字符串。
	ErrDurationNotString = errors.New("duration must be a JSON string")
	// ErrInvalidDuration 表示 Duration 字符串不符合 time.ParseDuration 格式。
	ErrInvalidDuration = errors.New("invalid duration")
)

// Duration 是 time.Duration 的包装，支持从 JSON 字符串解析人类可读的时间(如 "30s")。
type Duration time.Duration

// Std 把 task.Duration 转成标准的 time.Duration。
func (d Duration) Std() time.Duration {
	return time.Duration(d)
}

// UnmarshalJSON 让 Duration 支持 "30s" / "5m" 这类人类可读格式。
func (d *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("%w: %w", ErrDurationNotString, err)
	}

	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("%w %q: %w", ErrInvalidDuration, s, err)
	}

	*d = Duration(parsed)
	return nil
}
