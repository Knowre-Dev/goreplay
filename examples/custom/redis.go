package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"time"
)

type RedisMap struct {
	sync.Mutex
	field string
	rdb   *redis.Client
	ctx   context.Context
	m     map[string]int64
}

type AppConfig struct {
	MaxTTL int64
	Conn   *redis.Client
}

func NewRedisMap(c AppConfig, field string) *RedisMap {
	rdb := c.Conn
	m := make(map[string]int64)

	r := &RedisMap{
		rdb:   rdb,
		ctx:   context.Background(),
		m:     m,
		field: field,
	}

	go func() {
		for now := range time.Tick(time.Second) {
			r.Lock()
			for k, v := range r.m {
				if now.Unix()-v > int64(c.MaxTTL) {
					delete(r.m, k)
					r.rdb.Del(r.ctx, k)
				}
			}
			r.Unlock()
		}
	}()

	return r
}

func (r *RedisMap) Get(key string) (v string, exist bool) {
	val, err := r.rdb.HGet(r.ctx, key, r.field).Result()

	if err == redis.Nil {
		return "", false
	} else if err != nil {
		log.Fatal(err)
	}

	return val, true
}

func (r *RedisMap) Set(key string, value string) {
	err := r.rdb.HSet(r.ctx, key, r.field, value).Err()
	if err != nil {
		log.Fatal(err)
	}
	r.Lock()
	defer r.Unlock()
	r.m[key] = time.Now().Unix()
}

func (r *RedisMap) Len() int {
	r.Lock()
	defer r.Unlock()
	return len(r.m)
}

func (r *RedisMap) Del(key string) {
	err := r.rdb.HDel(r.ctx, key, r.field).Err()
	if err == redis.Nil {

	} else if err != nil {
		log.Fatal(err)
	}
}
