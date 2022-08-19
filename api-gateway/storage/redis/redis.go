package redis

import (
	redis "github.com/gomodule/redigo/redis"
	"github.com/my_tdm/api-gateway/storage/repo"
)

type redisRepo struct {
	rConn *redis.Pool
}

func NewRedisRepo(rds *redis.Pool) repo.RedisRepositoryStorage {
	return &redisRepo{
		rConn: rds,
	}
}

func (r *redisRepo) Set(key, value string) (err error) {
	conn := r.rConn.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	return
}

// SetWithTTL ...
func (r *redisRepo) SetWithTTL(key, value string, seconds int64) (err error) {
	conn := r.rConn.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, seconds, value)
	return
}

func (r *redisRepo) Get(key string) (interface{}, error) {
	conn := r.rConn.Get()
	defer conn.Close()

	return conn.Do("GET", key)
}
