package tools

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

var (
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
)

func TestRedisLock_TryLock(t *testing.T) {
	rl, err := NewRedisLock(client, "test-lock")
	if err != nil {
		t.Fatalf("Failed to create distribute lock: %v", err)
	}

	// Test TryLock() with valid ctx and key
	_, err = rl.TryLock(context.Background())
	if err != nil {
		t.Fatalf("Failed to lock: %v", err)
	}

	// Test TryLock() with invalid ctx
	_, err = rl.TryLock(context.TODO())
	if err == nil {
		t.Fatalf("Failed to lock with invalid ctx")
	}

	// Test TryLock() with invalid key
	rl.Client.SetNX(context.Background(), "", rl.uuid, defaultExp).Result()
	_, err = rl.TryLock(context.TODO())
	if err == nil {
		t.Fatalf("Failed to lock with invalid key")
	}
}

//func TestRedisLock_SpinLock(t *test.T) {
//	client := test.NewMachineClient()
//	key := "test-lock"
//	rl := NewRedisLock(client, key)
//
//	// Test SpinLock() with valid ctx and key
//	for i := 0; i < 5; i++ {
//		_, err := rl.SpinLock(context.Background(), 3)
//		if err != nil {
//			t.Fatalf("Failed to lock: %v", err)
//		}
//		time.Sleep(100 * time.Millisecond)
//	}
//
//	// Test SpinLock() with invalid ctx
//	for i := 0; i < 5; i++ {
//		_, err := rl.SpinLock(context.TODO(), 3)
//		if err == nil {
//			t.Fatalf("Failed to lock with invalid ctx")
//		}
//		time.Sleep(100 * time.Millisecond)
//	}
//
//	// Test SpinLock() with invalid key
//	rl.Client.SetNX(context.Background(), key, uuid.NewV4(), defaultExp).Result()
//	for i := 0; i < 5; i++ {
//		_, err := rl.SpinLock(context.TODO(), 3)
//		if err == nil {
//			t.Fatalf("Failed to lock with invalid key")
//		}
//		time.Sleep(100 * time.Millisecond)
//	}
//}
