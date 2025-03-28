package lock

import "context"

type Locker interface {
	Lock(ctx context.Context, key string) (Unlocker, error)
	TryLock(ctx context.Context, key string) (Unlocker, error)
	Unlock(ctx context.Context, key string) error
}

type Unlocker interface {
	Unlock(ctx context.Context) error
}
