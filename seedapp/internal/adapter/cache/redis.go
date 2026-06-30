package cache

import (
	"context"
	"seedapp/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
	"github.com/runsystemid/gocache"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/cache.go -package=mocks
type RedisCache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string, value int64) (int64, error)
	Decrement(ctx context.Context, key string, value int64) (int64, error)
	Delete(ctx context.Context, key ...string) (int64, error)
	Get(ctx context.Context, key string, value interface{}) error
	Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	GetFiberStorage() fiber.Storage
}
type Redis struct {
	Config       *config.Config `inject:"config"`
	Client       *gocache.Redis
	FiberStorage fiber.Storage
}

func (r *Redis) Startup() error {
	r.Client = gocache.New(gocache.RedisConfig{
		Mode:     r.Config.Redis.Mode,
		Address:  r.Config.Redis.Address,
		Port:     r.Config.Redis.Port,
		Password: r.Config.Redis.Password,
	})

	r.FiberStorage = redis.New(redis.Config{
		Host:     r.Config.Redis.Address,
		Password: r.Config.Redis.Password,
		Port:     r.Config.Redis.Port,
	})

	return nil
}

func (r *Redis) Shutdown() error {
	return r.Client.Close()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	return r.Client.Exists(ctx, key)
}

func (r *Redis) Increment(ctx context.Context, key string, value int64) (int64, error) {
	return r.Client.Increment(ctx, key, value)
}

func (r *Redis) Decrement(ctx context.Context, key string, value int64) (int64, error) {
	return r.Client.Decrement(ctx, key, value)
}

func (r *Redis) Delete(ctx context.Context, key ...string) (int64, error) {
	return r.Client.Delete(ctx, key...)
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	return r.Client.Get(ctx, key, value)
}

func (r *Redis) Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Put(ctx, key, value, expiration)
}

func (r *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.Client.Client.SetNX(ctx, key, value, expiration).Result()
}

func (r *Redis) GetFiberStorage() fiber.Storage {
	return r.FiberStorage
}
