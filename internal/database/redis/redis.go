package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"project-go-/internal/config"
)

type Context struct {
	config  *config.Config
	rClient *redis.Client
}

func CreateNewRedisContext(cfg *config.Config) (*Context, error) {
	log.Println("Connecting to Redis...")
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Address, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	log.Printf("Connected to Redis at %s", cfg.Redis.Address)
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
