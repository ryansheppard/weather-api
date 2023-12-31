package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Client *redis.Client
	Ctx    context.Context
}

func New(ctx context.Context, address string, database int) *Cache {
	if address == "" {
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       database,
	})

	return &Cache{
		Client: client,
		Ctx:    ctx,
	}
}

func (c *Cache) SetKey(key string, value interface{}, ttl int) error {
	seconds := time.Duration(ttl) * time.Second

	return c.Client.Set(c.Ctx, key, value, seconds).Err()
}

func (c *Cache) GetKey(key string) (interface{}, error) {
	result, err := c.Client.Get(c.Ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Cache) DeleteKey(key string) {
	// if cache == nil {
	// 	return
	// }

	c.Client.Del(c.Ctx, key)
}
