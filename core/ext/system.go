package ext

import "context"

type System interface {
    SystemInfo(ctx context.Context) (map[string]any, error)
}
