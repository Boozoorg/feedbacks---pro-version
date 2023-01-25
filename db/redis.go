package db

import (
	"log"

	"github.com/go-redis/redis"
)

var Rdb *redis.Client

func setupRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	log.Println("Success Redis connection")
}

func closeRedis() {
	err := Rdb.Close()
	if err != nil {
		log.Println("failed to close redis")
		return
	}
	log.Println("Redis closed")
}
