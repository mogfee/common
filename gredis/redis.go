package gredis

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var redisClient *redis.Client

func New(host, password string, db int) error {
	redisClient = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     host,
		Password: password,
		DB:       db,
	})

	return redisClient.Ping().Err()
}
func Client() *redis.Client {
	return redisClient
}

func Shutdown() {
	redisClient.ShutdownSave()
}

type lock struct {
	keys []string
}

func NewLock(prefix string, key ...string) *lock {
	keys := []string{}
	if len(key) > 0 {
		for _, v := range key {
			if v != "" {
				keys = append(keys, fmt.Sprintf("%s%s", prefix, v))
			}
		}
	}
	return &lock{
		keys: keys,
	}
}
func (l *lock) Lock() bool {
	oks := []string{}
	for _, v := range l.keys {
		ok, _ := redisClient.SetNX(v, 1, time.Second*5).Result()
		log.Printf("add lock :%s ok:%v\n", v, ok)
		if ok {
			oks = append(oks, v)
		} else {
			//加锁失败
			for _, v := range oks {
				log.Printf("error delete lock :%s\n", v)
				redisClient.Del(v)
			}
			return false
		}
	}
	return true
}
func (l *lock) UnLock() {
	for _, v := range l.keys {
		log.Printf("delete lock :%s\n", v)
		redisClient.Del(v)
	}
}
