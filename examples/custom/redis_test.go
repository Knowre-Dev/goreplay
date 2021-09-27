package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisMapGetSet(t *testing.T) {
	conn := redis.NewClient(&redis.Options{
		Addr:     "debtimes.com:27017",
		Password: "",
		DB:       0,
	})

	appConfig := AppConfig{
		MaxTTL: 10,
		Conn:   conn,
	}
	xToken := NewRedisMap(appConfig, "x-access-token")
	xToken.Set("element", "abcdef")
	if value, exist := xToken.Get("element"); exist {
		assert.Equal(t, value, "abcdef")
	}

	cookie := NewRedisMap(appConfig, "set-cookie")
	cookie.Set("element", "coooooookie")
	if value, exist := cookie.Get("element"); exist {
		assert.Equal(t, value, "coooooookie")
	}

	xToken.Set("element", "bbbbbb")
	if value, exist := xToken.Get("element"); exist {
		assert.Equal(t, value, "bbbbbb")
	}

	time.Sleep(time.Second * time.Duration(appConfig.MaxTTL+2))
	_, exist := xToken.Get("element")
	assert.False(t, exist)

}
