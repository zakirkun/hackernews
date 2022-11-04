package redis

import (
	db "github.com/go-redis/redis"
)

type redisContext struct{}

func New() *redisContext {
	return &redisContext{}
}

func (r *redisContext) Set(key string, value interface{}) *error {
	con := client()
	defer con.Close()

	err := con.Set(key, value, 0).Err()
	if err != nil {
		return &err
	}

	return nil
}

func (r *redisContext) Get(key string) (string, *error) {
	con := client()
	defer con.Close()

	val, err := con.Get(key).Result()
	if err != nil {
		return "", &err
	}

	return val, nil
}

func client() *db.Client {
	rdb := db.NewClient(&db.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
