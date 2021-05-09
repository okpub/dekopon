package utils

import "context"

var (
	ParentCtxKey = struct{}{}
	emptCtxValue = struct{}{}
)

//仅仅包裹一层
func WithValue(ctx context.Context) context.Context {
	return context.WithValue(ctx, ParentCtxKey, emptCtxValue)
}
