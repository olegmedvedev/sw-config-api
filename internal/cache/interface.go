package cache

import "time"

// Interface defines the cache interface
type Interface interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, ttl time.Duration) error
	Delete(key string) error
	Close() error
}
