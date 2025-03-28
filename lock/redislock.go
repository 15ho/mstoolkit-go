package lock

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrLockFailed = errors.New("lock failed")
)

type redisLockerConfig struct {
	lockTTL         time.Duration
	keyPrefix       string
	maxAcquireCount int
	acquireInterval time.Duration
}

type RedisLockerOption interface {
	apply(opt *redisLockerConfig)
}
type redisLockerOptionFunc func(opt *redisLockerConfig)

func (f redisLockerOptionFunc) apply(opt *redisLockerConfig) {
	f(opt)
}

func WithRedisLockerLockTTL(ttl time.Duration) RedisLockerOption {
	return redisLockerOptionFunc(func(opt *redisLockerConfig) {
		opt.lockTTL = ttl
	})
}

func WithRedisLockerKeyPrefix(prefix string) RedisLockerOption {
	return redisLockerOptionFunc(func(opt *redisLockerConfig) {
		opt.keyPrefix = prefix
	})
}
func WithRedisLockerMaxAcquireCount(count int) RedisLockerOption {
	return redisLockerOptionFunc(func(opt *redisLockerConfig) {
		opt.maxAcquireCount = count
	})
}
func WithRedisLockerAcquireInterval(interval time.Duration) RedisLockerOption {
	return redisLockerOptionFunc(func(opt *redisLockerConfig) {
		opt.acquireInterval = interval
	})
}

type redisLocker struct {
	*redisUnlocker
	rd  *redis.Client
	cfg *redisLockerConfig
}

func NewRedisLocker(rd *redis.Client, opts ...RedisLockerOption) (Locker, error) {
	if rd == nil {
		return nil, errors.New("redis client is nil")
	}
	cfg := &redisLockerConfig{}
	for _, opt := range opts {
		opt.apply(cfg)
	}
	if cfg.lockTTL == 0 {
		cfg.lockTTL = 5 * time.Second
	}
	if cfg.keyPrefix == "" {
		cfg.keyPrefix = "redislock:"
	}
	if cfg.maxAcquireCount == 0 {
		cfg.maxAcquireCount = 3
	}
	if cfg.acquireInterval == 0 {
		cfg.acquireInterval = 100 * time.Millisecond
	}
	return &redisLocker{redisUnlocker: &redisUnlocker{rd: rd}, rd: rd, cfg: cfg}, nil
}

func (r *redisLocker) Lock(ctx context.Context, key string) (Unlocker, error) {
	key = r.cfg.keyPrefix + key
	for range r.cfg.maxAcquireCount {
		ok, err := r.rd.SetNX(ctx, key, 1, r.cfg.lockTTL).Result()
		if err != nil {
			return nil, err
		}
		if ok {
			return &redisUnlocker{rd: r.rd, key: key}, nil
		}
		time.Sleep(r.cfg.acquireInterval)
	}
	return nil, ErrLockFailed
}

func (r *redisLocker) TryLock(ctx context.Context, key string) (Unlocker, error) {
	key = r.cfg.keyPrefix + key
	ok, err := r.rd.SetNX(ctx, key, 1, r.cfg.lockTTL).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrLockFailed
	}
	return &redisUnlocker{rd: r.rd, key: key}, nil
}

func (r *redisLocker) Unlock(ctx context.Context, key string) error {
	r.redisUnlocker.key = r.cfg.keyPrefix + key
	return r.redisUnlocker.Unlock(ctx)
}

type redisUnlocker struct {
	rd  *redis.Client
	key string
}

func (u *redisUnlocker) Unlock(ctx context.Context) error {
	_, err := u.rd.Del(ctx, u.key).Result()
	if err != nil {
		return err
	}
	return nil
}
