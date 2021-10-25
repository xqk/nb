package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	defaultDepth = 3
	// DefaultCaller 返回文件和行的值
	DefaultCaller = Caller(defaultDepth)

	// DefaultTimestamp 返回当前时钟的值
	DefaultTimestamp = Timestamp(time.RFC3339)
)

// Valuer 返回一个日志的值
type Valuer func(ctx context.Context) interface{}

// Value 返回方法的值
func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

// Caller 返回一个pkg/file:line描述的值
func Caller(depth int) Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		if strings.LastIndex(file, "/log/filter.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}
		if strings.LastIndex(file, "/log/helper.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// Timestamp 返回自定义时间格式的时间戳的值
func Timestamp(layout string) Valuer {
	return func(context.Context) interface{} {
		return time.Now().Format(layout)
	}
}

func bindValues(ctx context.Context, keyvals []interface{}) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v(ctx)
		}
	}
}

func containsValuer(keyvals []interface{}) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}
