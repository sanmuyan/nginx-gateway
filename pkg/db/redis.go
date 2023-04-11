package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"nginx-gateway/pkg/config"
)

var RDB *redis.Client

func InitRedis() {
	var ctx = context.Background()
	RDB = redis.NewClient(&redis.Options{
		Addr: config.Conf.RedisAddr,
	})
	if ok, err := RDB.Ping(ctx).Result(); ok != "PONG" && err != nil {
		log.Println(ok)
		panic(err)
	}

}
