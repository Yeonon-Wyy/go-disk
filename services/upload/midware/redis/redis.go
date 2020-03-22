package redis

import (
	"github.com/go-redis/redis/v7"
	"go-disk/services/upload/config"
)

var (
	DSConfig = config.Conf.DataSource
)

func FSConn() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     DSConfig.Redis.Addr,
		Password: DSConfig.Redis.Password,
		DB:       DSConfig.Redis.Database,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
