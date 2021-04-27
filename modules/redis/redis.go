package redis

import "github.com/go-redis/redis"

type DB struct {
	d *redis.Client
}

// NewDB return abstract redis database connection
func NewDB(addr string, pwd string) *DB {
	return &DB{
		d: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pwd,
		}),
	}
}
