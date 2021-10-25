package metadata

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		mds []map[string]string
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "hello",
			args: args{[]map[string]string{{"hello": "nb"}, {"hello2": "go-nb"}}},
			want: Metadata{"hello": "nb", "hello2": "go-nb"},
		},
		{
			name: "hi",
			args: args{[]map[string]string{{"hi": "nb"}, {"hi2": "go-nb"}}},
			want: Metadata{"hi": "nb", "hi2": "go-nb"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.mds...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want string
	}{
		{
			name: "nb",
			m:    Metadata{"nb": "value", "env": "dev"},
			args: args{key: "nb"},
			want: "value",
		},
		{
			name: "env",
			m:    Metadata{"nb": "value", "env": "dev"},
			args: args{key: "env"},
			want: "dev",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetadata_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		m    Metadata
		args args
		want Metadata
	}{
		{
			name: "nb",
			m:    Metadata{},
			args: args{key: "hello", value: "nb"},
			want: Metadata{"hello": "nb"},
		},
		{
			name: "env",
			m:    Metadata{"hello": "nb"},
			args: args{key: "env", value: "pro"},
			want: Metadata{"hello": "nb", "env": "pro"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Set(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Set() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestClientContext(t *testing.T) {
	type args struct {
		ctx context.Context
		md  Metadata
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "nb",
			args: args{context.Background(), Metadata{"hello": "nb", "nb": "https://www.icl.site"}},
		},
		{
			name: "hello",
			args: args{context.Background(), Metadata{"hello": "nb", "hello2": "https://www.icl.site"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewClientContext(tt.args.ctx, tt.args.md)
			m, ok := FromClientContext(ctx)
			if !ok {
				t.Errorf("FromClientContext() = %v, want %v", ok, true)
			}

			if !reflect.DeepEqual(m, tt.args.md) {
				t.Errorf("meta = %v, want %v", m, tt.args.md)
			}
		})
	}
}

func TestServerContext(t *testing.T) {
	type args struct {
		ctx context.Context
		md  Metadata
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "nb",
			args: args{context.Background(), Metadata{"hello": "nb", "nb": "https://www.icl.site"}},
		},
		{
			name: "hello",
			args: args{context.Background(), Metadata{"hello": "nb", "hello2": "https://www.icl.site"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewServerContext(tt.args.ctx, tt.args.md)
			m, ok := FromServerContext(ctx)
			if !ok {
				t.Errorf("FromServerContext() = %v, want %v", ok, true)
			}

			if !reflect.DeepEqual(m, tt.args.md) {
				t.Errorf("meta = %v, want %v", m, tt.args.md)
			}
		})
	}
}

func TestAppendToClientContext(t *testing.T) {
	type args struct {
		md Metadata
		kv []string
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "nb",
			args: args{Metadata{}, []string{"hello", "nb", "env", "dev"}},
			want: Metadata{"hello": "nb", "env": "dev"},
		},
		{
			name: "hello",
			args: args{Metadata{"hi": "https://www.icl.site"}, []string{"hello", "nb", "env", "dev"}},
			want: Metadata{"hello": "nb", "env": "dev", "hi": "https://www.icl.site"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewClientContext(context.Background(), tt.args.md)
			ctx = AppendToClientContext(ctx, tt.args.kv...)
			md, ok := FromClientContext(ctx)
			if !ok {
				t.Errorf("FromServerContext() = %v, want %v", ok, true)
			}
			if !reflect.DeepEqual(md, tt.want) {
				t.Errorf("metadata = %v, want %v", md, tt.want)
			}
		})
	}
}

func TestMergeToClientContext(t *testing.T) {
	type args struct {
		md       Metadata
		appendMd Metadata
	}
	tests := []struct {
		name string
		args args
		want Metadata
	}{
		{
			name: "nb",
			args: args{Metadata{}, Metadata{"hello": "nb", "env": "dev"}},
			want: Metadata{"hello": "nb", "env": "dev"},
		},
		{
			name: "hello",
			args: args{Metadata{"hi": "https://www.icl.site"}, Metadata{"hello": "nb", "env": "dev"}},
			want: Metadata{"hello": "nb", "env": "dev", "hi": "https://www.icl.site"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewClientContext(context.Background(), tt.args.md)
			ctx = MergeToClientContext(ctx, tt.args.appendMd)
			md, ok := FromClientContext(ctx)
			if !ok {
				t.Errorf("FromServerContext() = %v, want %v", ok, true)
			}
			if !reflect.DeepEqual(md, tt.want) {
				t.Errorf("metadata = %v, want %v", md, tt.want)
			}
		})
	}
}

func TestMetadata_Range(t *testing.T) {
	md := Metadata{"nb": "nb", "https://www.icl.site": "https://www.icl.site", "go-nb": "go-nb"}
	tmp := Metadata{}
	md.Range(func(k, v string) bool {
		if k == "https://www.icl.site" || k == "nb" {
			tmp[k] = v
		}
		return true
	})
	if !reflect.DeepEqual(tmp, Metadata{"https://www.icl.site": "https://www.icl.site", "nb": "nb"}) {
		t.Errorf("metadata = %v, want %v", tmp, Metadata{"nb": "nb"})
	}
}

func TestMetadata_Clone(t *testing.T) {
	tests := []struct {
		name string
		m    Metadata
		want Metadata
	}{
		{
			name: "nb",
			m:    Metadata{"nb": "nb", "https://www.icl.site": "https://www.icl.site", "go-nb": "go-nb"},
			want: Metadata{"nb": "nb", "https://www.icl.site": "https://www.icl.site", "go-nb": "go-nb"},
		},
		{
			name: "go",
			m:    Metadata{"language": "golang"},
			want: Metadata{"language": "golang"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Clone()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
			got["nb"] = "go"
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("want got != want got %v want %v", got, tt.want)
			}
		})
	}
}
