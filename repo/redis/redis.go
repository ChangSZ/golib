package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var _ Repo = (*cacheRepo)(nil)

type Repo interface {
	i()
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Expire(ctx context.Context, key string, ttl time.Duration) bool
	ExpireAt(ctx context.Context, key string, ttl time.Time) bool
	Del(ctx context.Context, key string) bool
	Exists(ctx context.Context, keys ...string) bool
	Incr(ctx context.Context, key string) int64
	Close() error
	Version() string
}

var cache *cacheRepo

type cacheRepo struct {
	client *redis.Client
}

func Init(cfg Config) {
	client, err := redisConnect(cfg)
	if err != nil {
		panic(fmt.Sprintf("redis连接失败: %v", err))
	}

	// 设置钩子函数, 个人觉着没必要, 如有需要自己实现下hook
	// client.AddHook(nil)

	cache = &cacheRepo{client: client}
}

type Config struct {
	Addr         string `toml:"addr"`
	Pass         string `toml:"pass"`
	Db           int    `toml:"db"`
	MaxRetries   int    `toml:"maxRetries"`
	PoolSize     int    `toml:"poolSize"`
	MinIdleConns int    `toml:"minIdleConns"`
}

func Cache() Repo {
	return cache
}

func (c *cacheRepo) i() {}

func redisConnect(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pass,
		DB:           cfg.Db,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	if err := client.Ping(client.Context()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis err: %w", err)
	}

	return client, nil
}

func (c *cacheRepo) Client() *redis.Client {
	return c.client
}

// Set set some <key,value> into redis
func (c *cacheRepo) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("redis set key: %s err: %w", key, err)
	}

	return nil
}

func (c *cacheRepo) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// TTL get some key from redis
func (c *cacheRepo) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return -1, fmt.Errorf("redis get key: %s err: %w", key, err)
	}
	return ttl, nil
}

// Expire expire some key
func (c *cacheRepo) Expire(ctx context.Context, key string, ttl time.Duration) bool {
	ok, _ := c.client.Expire(ctx, key, ttl).Result()
	return ok
}

// ExpireAt expire some key at some time
func (c *cacheRepo) ExpireAt(ctx context.Context, key string, ttl time.Time) bool {
	ok, _ := c.client.ExpireAt(ctx, key, ttl).Result()
	return ok
}

func (c *cacheRepo) Exists(ctx context.Context, keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	value, _ := c.client.Exists(ctx, keys...).Result()
	return value > 0
}

func (c *cacheRepo) Del(ctx context.Context, key string) bool {
	if key == "" {
		return true
	}
	value, _ := c.client.Del(ctx, key).Result()
	return value > 0
}

func (c *cacheRepo) Incr(ctx context.Context, key string) int64 {
	value, _ := c.client.Incr(ctx, key).Result()
	return value
}

// Close close redis client
func (c *cacheRepo) Close() error {
	return c.client.Close()
}

// Version redis server version
func (c *cacheRepo) Version() string {
	server := c.client.Info(context.Background(), "server").Val()
	spl1 := strings.Split(server, "# Server")
	spl2 := strings.Split(spl1[1], "redis_version:")
	spl3 := strings.Split(spl2[1], "redis_git_sha1:")
	return spl3[0]
}
