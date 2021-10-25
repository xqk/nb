package log

// FilterOption 过滤器选项
type FilterOption func(*Filter)

const fuzzyStr = "***"

// FilterLevel 过滤器等级
func FilterLevel(level Level) FilterOption {
	return func(opts *Filter) {
		opts.level = level
	}
}

// FilterKey 添加过滤键
func FilterKey(key ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range key {
			o.key[v] = struct{}{}
		}
	}
}

// FilterValue 添加过滤值
func FilterValue(value ...string) FilterOption {
	return func(o *Filter) {
		for _, v := range value {
			o.value[v] = struct{}{}
		}
	}
}

// FilterFunc 添加过滤方法
func FilterFunc(f func(level Level, keyvals ...interface{}) bool) FilterOption {
	return func(o *Filter) {
		o.filter = f
	}
}

// Filter 日志过滤器
type Filter struct {
	logger Logger
	level  Level
	key    map[interface{}]struct{}
	value  map[interface{}]struct{}
	filter func(level Level, keyvals ...interface{}) bool
}

// NewFilter 创建一个日志过滤器
func NewFilter(logger Logger, opts ...FilterOption) *Filter {
	options := Filter{
		logger: logger,
		key:    make(map[interface{}]struct{}),
		value:  make(map[interface{}]struct{}),
	}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

// Log 根据等级和键值打印日志
func (f *Filter) Log(level Level, keyvals ...interface{}) error {
	if level < f.level {
		return nil
	}
	if f.filter != nil && f.filter(level, keyvals...) {
		return nil
	}
	if len(f.key) > 0 || len(f.value) > 0 {
		for i := 0; i < len(keyvals); i += 2 {
			v := i + 1
			if v >= len(keyvals) {
				continue
			}
			if _, ok := f.key[keyvals[i]]; ok {
				keyvals[v] = fuzzyStr
			}
			if _, ok := f.value[keyvals[v]]; ok {
				keyvals[v] = fuzzyStr
			}
		}
	}
	return f.logger.Log(level, keyvals...)
}
