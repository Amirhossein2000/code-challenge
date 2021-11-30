package main

import (
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"sync"
)

var (
	redisClient   *redis.Client
	redisDB       = 2
	redisChan     = "orders"
	redisAddress  string
	redisPassword string
)

func init() {
	strDB := os.Getenv("REDIS_DB")
	if strDB != "" {
		db, err := strconv.Atoi(strDB)
		if err != nil {
			logger.Warningf("REDIS_DB convert err: %s", err.Error())
		} else {
			redisDB = db
		}
	}

	rc := os.Getenv("REDIS_DB")
	if rc != "" {
		redisChan = rc
	}

	redisAddress = os.Getenv("REDIS_ADDRESS")
	redisPassword = os.Getenv("REDIS_PASSWORD")
}

func getRedisClient() *redis.Client {
	var doOnce sync.Once
	doOnce.Do(
		func() {
			redisClient = redis.NewClient(&redis.Options{
				Addr:     redisAddress,
				Password: redisPassword,
				DB:       redisDB,
			})
			err := redisClient.Ping().Err()
			if err != nil {
				logger.Errorf("redis ping err: %s", err.Error())
			}
		})
	return redisClient
}
