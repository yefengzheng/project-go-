package database

import (
	"project-go-/internal/config"
	"project-go-/internal/database/redis"
	"time"
)

// for our db execution context
type Context struct {
	config       *config.Config
	RedisContext *redis.Context
	PgsqlContext *pgsql.Context
}

func CreateNewDbContext(config *config.Config, lifeTime time.Duration) (dbContext *Context, err error) {
	rCtx, err := redis.CreateNewRedisContext(config) //if()
	dbContext.RedisContext = rCtx
	return dbContext, nil
}
