package redis

import (
	"github.com/go-redis/redis/v7"
	"go-disk/config"
)

func AuthConn() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: config.AuthRedisAddr,
		Password: config.AuthRedisPassword,
		DB: config.AuthRedisDB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func FSConn() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: config.FSRedisAddr,
		Password: config.FSRedisPassword,
		DB: config.FSRedisDB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

