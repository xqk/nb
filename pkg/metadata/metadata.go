package metadata

import (
	"context"
	"fmt"
	"strings"
)

// Metadata 微服务之间将会通过 HTTP 和 gRPC 进行接口交互，所以服务架构需要统一元信息传递使用。
// 目前 gRPC 中也可以携带元信息传递，原理是放到 HTTP Header 中，然后上游将会收到对应的元信息信息。
// 所以设计上，也是通过 HTTP Header 进行传递，在框架中先通过元信息包封装成key/value结构，然后携带到 Transport Header 中。
type Metadata map[string]string

// New 根据给定的键值对map创建一个元信息
func New(mds ...map[string]string) Metadata {
	md := Metadata{}
	for _, m := range mds {
		for k, v := range m {
			md.Set(k, v)
		}
	}
	return md
}

// Get 返回与传递的键相关联的值
func (m Metadata) Get(key string) string {
	k := strings.ToLower(key)
	return m[k]
}

// Set 存储键值对
func (m Metadata) Set(key string, value string) {
	if key == "" || value == "" {
		return
	}
	k := strings.ToLower(key)
	m[k] = value
}

// Range 在元信息中迭代元素
func (m Metadata) Range(f func(k, v string) bool) {
	for k, v := range m {
		ret := f(k, v)
		if !ret {
			break
		}
	}
}

// Clone 返回一个深度拷贝的元信息
func (m Metadata) Clone() Metadata {
	md := Metadata{}
	for k, v := range m {
		md[k] = v
	}
	return md
}

type serverMetadataKey struct{}

// NewServerContext 创建一个附加服务端元信息的上下文
func NewServerContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, serverMetadataKey{}, md)
}

// FromServerContext 如果服务端的上下文存在元信息则返回元信息
func FromServerContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(serverMetadataKey{}).(Metadata)
	return md, ok
}

type clientMetadataKey struct{}

// NewClientContext 创建一个附加客户端元信息的上下文
func NewClientContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, clientMetadataKey{}, md)
}

// FromClientContext 如果客户端上下文中存在元信息则返回元信息
func FromClientContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(clientMetadataKey{}).(Metadata)
	return md, ok
}

// AppendToClientContext 返回一个新上下文，其中所提供的kv与上下文中的任何现有元信息合并。
func AppendToClientContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToOutgoingContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for i := 0; i < len(kv); i += 2 {
		md.Set(kv[i], kv[i+1])
	}
	return NewClientContext(ctx, md)
}

// MergeToClientContext 合并一个新的元信息到上下文
func MergeToClientContext(ctx context.Context, cmd Metadata) context.Context {
	md, _ := FromClientContext(ctx)
	md = md.Clone()
	for k, v := range cmd {
		md[k] = v
	}
	return NewClientContext(ctx, md)
}
