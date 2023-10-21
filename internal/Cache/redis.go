package Cache

import (
	"github.com/go-redis/redis/v8"
	cache_interface "short-link/internal/Cache/Interface"
)

// Db holds database connection to Postgres
type Cache struct {
	client *redis.Client
}

func Connect() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
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

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateCache() cache_interface.CacheInterface {

	client := Connect()

	cache := &Cache{client}

	return cache
}
