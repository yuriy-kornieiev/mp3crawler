package main

import (
	"github.com/go-redis/redis"
	"log"
)

type Redis struct {
	Host string
	Port string

	Connection *redis.Client
}

func (r Redis) Load() map[string]*Redis {

	data := make(map[string]*Redis)

	data["yuriyk"] = &Redis{
		"127.0.0.1",
		"6379",
		nil,
	}

	return data

}

func (r *Redis) GetConnection() *redis.Client {

	if r.Connection != nil {
		return r.Connection
	}

	client := redis.NewClient(&redis.Options{
		Addr:     r.Host + ":" + r.Port,
		Password: "",
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Panic(err)
	}

	r.Connection = client
	return client
}
