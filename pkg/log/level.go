package log

import "strings"

// Level 日志的等级
type Level int8

// LevelKey 日志等级的键
const LevelKey = "level"

const (
	// LevelDebug 日志debug等级
	LevelDebug Level = iota - 1
	// LevelInfo 日志info等级
	LevelInfo
	// LevelWarn 日志warn等级
	LevelWarn
	// LevelError 日志error等级
	LevelError
	// LevelFatal 日志fatal等级
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// ParseLevel 将等级字符串映射为日志的等级
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "FATAL":
		return LevelFatal
	}
	return LevelInfo
}
