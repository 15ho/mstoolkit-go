package lock

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisLocker(t *testing.T) {
	s := assert.New(t)

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if !s.NoError(client.FlushAll(context.Background()).Err()) {
		return
	}
	defer client.Close()

	locker, err := NewRedisLocker(client,
		WithRedisLockerAcquireInterval(200*time.Millisecond),
		WithRedisLockerMaxAcquireCount(5),
		WithRedisLockerKeyPrefix("test:"),
		WithRedisLockerLockTTL(10*time.Second))

	if !s.NoError(err) {
		return
	}

	ch := make(chan Unlocker)

	go func() {
		unlocker, err := locker.Lock(context.Background(), "key1")
		if !s.NoError(err) {
			close(ch)
			return
		}
		ch <- unlocker
		_, err = locker.Lock(context.Background(), "key2")
		s.NoError(err)
	}()

	unlocker, ok := <-ch
	if !s.True(ok) {
		return
	}

	if !s.NoError(unlocker.Unlock(context.Background())) {
		return
	}

	s.NoError(locker.Unlock(context.Background(), "key2"))
}
