package redis

import (
	"time"

	"github.com/go-redis/redis"
	"project-go-/internal/config"
)

type Context struct {
	config  *config.Config
	rClient *redis.Client
}

func CreateNewRedisContext(cfg *config.Config) (*Context, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &Context{
		config:  cfg,
		rClient: client,
	}, nil
}

func (r *Context) CheckHealth() (string, error) {
	return r.rClient.Ping().Result()
}

func (r *Context) Cleanup() {
	_ = r.rClient.Close()
}

func (r *Context) SetKeyValue(key, val string, duration time.Duration) error {
	return r.rClient.Set(key, val, duration).Err()
}

func (r *Context) DeleteKey(key string) {
	_ = r.rClient.Del(key).Err()
}

func (r *Context) GetValue(key string) (string, error) {
	return r.rClient.Get(key).Result()
}

func (r *Context) GetValidTime(key string) (time.Duration, error) {
	return r.rClient.TTL(key).Result()
}
