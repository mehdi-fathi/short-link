package Ports

import "time"

type MemCacheInterface interface {
	GetSlice(key string) ([]interface{}, bool)
	SetSlice(key string, content []interface{}, duration time.Duration)
}
