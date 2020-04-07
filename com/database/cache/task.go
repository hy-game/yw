package cache

import "github.com/garyburd/redigo/redis"

type IRedisTask interface {
	Run(redis.Conn)
}
