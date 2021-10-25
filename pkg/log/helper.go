package log

import (
	"context"
	"fmt"
	"os"
)

var DefaultMessageKey = "msg"

// Option 助手选项
type Option func(*Helper)

// Helper 一个日志的助手
type Helper struct {
	logger Logger
	msgKey string
}

func WithMessageKey(k string) Option {
	return func(opts *Helper) {
		opts.msgKey = k
	}
}

// NewHelper 创建一个日志的助手
func NewHelper(logger Logger, opts ...Option) *Helper {
	options := &Helper{
		msgKey: DefaultMessageKey, // default message key
		logger: logger,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithContext 返回h的一个浅拷贝，其上下文是ctx。提供的ctx必须是non-nil
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey: h.msgKey,
		logger: WithContext(ctx, h.logger),
	}
}

// Log 根据等级和键值对打印日志
func (h *Helper) Log(level Level, keyvals ...interface{}) {
	_ = h.logger.Log(level, keyvals...)
}

// Debug 记录debug级别的日志
func (h *Helper) Debug(a ...interface{}) {
	h.Log(LevelDebug, h.msgKey, fmt.Sprint(a...))
}

// Debugf 记录debug级别的日志
func (h *Helper) Debugf(format string, a ...interface{}) {
	h.Log(LevelDebug, h.msgKey, fmt.Sprintf(format, a...))
}

// Debugw 记录debug级别的日志
func (h *Helper) Debugw(keyvals ...interface{}) {
	h.Log(LevelDebug, keyvals...)
}

// Info 记录info级别的日志
func (h *Helper) Info(a ...interface{}) {
	h.Log(LevelInfo, h.msgKey, fmt.Sprint(a...))
}

// Infof 记录info级别的日志
func (h *Helper) Infof(format string, a ...interface{}) {
	h.Log(LevelInfo, h.msgKey, fmt.Sprintf(format, a...))
}

// Infow 记录info级别的日志
func (h *Helper) Infow(keyvals ...interface{}) {
	h.Log(LevelInfo, keyvals...)
}

// Warn 记录warn级别的日志
func (h *Helper) Warn(a ...interface{}) {
	h.Log(LevelWarn, h.msgKey, fmt.Sprint(a...))
}

// Warnf 记录warn级别的日志
func (h *Helper) Warnf(format string, a ...interface{}) {
	h.Log(LevelWarn, h.msgKey, fmt.Sprintf(format, a...))
}

// Warnw 记录warn级别的日志
func (h *Helper) Warnw(keyvals ...interface{}) {
	h.Log(LevelWarn, keyvals...)
}

// Error 记录error级别的日志
func (h *Helper) Error(a ...interface{}) {
	h.Log(LevelError, h.msgKey, fmt.Sprint(a...))
}

// Errorf 记录error级别的日志
func (h *Helper) Errorf(format string, a ...interface{}) {
	h.Log(LevelError, h.msgKey, fmt.Sprintf(format, a...))
}

// Errorw 记录error级别的日志
func (h *Helper) Errorw(keyvals ...interface{}) {
	h.Log(LevelError, keyvals...)
}

// Fatal 记录fatal级别的日志
func (h *Helper) Fatal(a ...interface{}) {
	h.Log(LevelFatal, h.msgKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf 记录fatal级别的日志
func (h *Helper) Fatalf(format string, a ...interface{}) {
	h.Log(LevelFatal, h.msgKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw 记录fatal级别的日志
func (h *Helper) Fatalw(keyvals ...interface{}) {
	h.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
