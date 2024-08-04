package rd

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb *redis.Client = nil

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-16633.c61.us-east-1-3.ec2.redns.redis-cloud.com:16633",
		Password: "R3KdGsECUaAjAFRWgpMiAwdaQYWpRlzk",
		DB:       0,
	})
	
	ping, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Printf("Не удалось подключиться к Redis: %v", err)
		return
    }
    log.Printf("Подключение успешно: %s\n", ping)
}

func GetRdb() *redis.Client{
	return rdb
}

func GetCtx() context.Context{
	return ctx
}

func CloseRdb() {
	rdb.Close()
}