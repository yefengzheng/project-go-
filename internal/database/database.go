package database

import (
	"project-go-/internal/config"
	"project-go-/internal/database/pgsql"
	"project-go-/internal/database/redis"
	"time"
)

type Context struct {
	Config       *config.Config
	RedisContext *redis.Context
	PgsqlContext *pgsql.Context
}

func CreateNewDbContext(cfg *config.Config, lifeTime time.Duration) (*Context, error) {
	rCtx, err := redis.CreateNewRedisContext(cfg)
	if err != nil {
		return nil, err
	}

	pgCtx, err := pgsql.CreateNewPgsqlContext(cfg, lifeTime)
	if err != nil {
		return nil, err
	}

	return &Context{
		Config:       cfg,
		RedisContext: rCtx,
		PgsqlContext: pgCtx,
	}, nil
}

func (ctx *Context) Ping() error {
	return ctx.PgsqlContext.DB.Ping()
}
