package tools

import (
	"context"
	"github.com/gofrs/uuid"
	"time"
)

const (
	// defaultExp default timeout for lock
	defaultExp = 10 * time.Second
	// sleepDur default sleep time for spin lock
	sleepDur = 10 * time.Millisecond
)

type RedisLock struct {
	Client     RedisClient
	Key        string // resources that need to be locked
	uuid       string // lock owner uuid
	cancelFunc context.CancelFunc
}

// NewRedisLock create a redis distribute lock
func NewRedisLock(client RedisClient, key string) (*RedisLock, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &RedisLock{
		Client: client,
		Key:    key,
		uuid:   id.String(),
	}, nil
}

// TryLock attempt to lock, return true if success, otherwise false
func (rl *RedisLock) TryLock(ctx context.Context) (bool, error) {
	res, err := rl.Client.SetNX(ctx, rl.Key, rl.uuid, defaultExp).Result()
	if err != nil || !res {
		return false, err
	}
	c, cancel := context.WithCancel(ctx)
	rl.cancelFunc = cancel
	rl.refresh(c)
	return res, nil
}

// SpinLock Loop `retryTimes` times to call TryLock()
func (rl *RedisLock) SpinLock(ctx context.Context, retryTimes int) (bool, error) {
	for i := 0; i < retryTimes; i++ {
		resp, err := rl.TryLock(ctx)
		if err != nil {
			return false, err
		}
		if resp {
			return resp, nil
		}
		time.Sleep(sleepDur)
	}
	return false, nil
}

// Unlock attempt to unlock , return true if success, otherwise false
func (rl *RedisLock) Unlock(ctx context.Context) (bool, error) {
	resp, err := NewTools(rl.Client).Cad(ctx, rl.Key, rl.uuid)
	if err != nil {
		return false, err
	}
	if resp {
		rl.cancelFunc()
	}
	return resp, nil
}

func (rl *RedisLock) refresh(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(defaultExp / 4)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				rl.Client.Expire(ctx, rl.Key, defaultExp)
			}
		}
	}()
}
