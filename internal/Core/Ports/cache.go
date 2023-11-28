package Ports

type CacheInterface interface {
	Ping() (string, error)
	Hset(key string, values ...interface{})
	Hget(key, field string) (string, error)
	IncrBy(key string, value int64) (int64, error)
	Get(key string) (string, error)
}
