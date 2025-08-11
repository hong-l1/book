package ratelimit

import "context"

type Limit interface {
	Limited(ctx context.Context, key string) (bool, error)
}
