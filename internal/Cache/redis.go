package Cache

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	cache_interface "short-link/internal/Cache/Interface"
	"short-link/internal/Config"
)

// Db holds database connection to Postgres
type Cache struct {
	client *redis.Client
	Config *Config.Config
}

func Connect(cfg *Config.Config) *redis.Client {

	url := fmt.Sprintf("%s:%d", cfg.Redis.HOST, cfg.Redis.PORT)

	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: cfg.Redis.PASSWORD,
		DB:       0,
	})

	return client
}

func (Cache *Cache) Ping() (string, error) {

	pong, err := Cache.client.Ping(Cache.client.Context()).Result()
	return pong, err
}

func (Cache *Cache) Hset(key string, values ...interface{}) {

	Cache.client.HSet(Cache.client.Context(), key, values)
}

func (Cache *Cache) Hget(key, field string) (string, error) {

	return Cache.client.HGet(Cache.client.Context(), key, field).Result()
}
func (Cache *Cache) Get(key string) (string, error) {

	return Cache.client.Get(Cache.client.Context(), key).Result()
}
func (Cache *Cache) IncrBy(key string, value int64) (int64, error) {

	return Cache.client.IncrBy(Cache.client.Context(), key, value).Result()
}

func CreateCache(cfg *Config.Config) cache_interface.CacheInterface {

	client := Connect(cfg)

	cache := &Cache{client, cfg}

	// Send a PING command to check the connection.
	_, err := cache.Ping()
	if err != nil {
		panic(errors.New("Redis: Can't connect to redis").(interface{}))
	}

	return cache
}
