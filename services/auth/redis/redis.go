package redis

import (
	"github.com/go-redis/redis/v7"
	"go-disk/services/auth/config"
	"log"
)

var (
	redisConfig = config.Conf.DataSource.Redis
)

func AuthConn() (*redis.Client, error) {
	log.Printf("[DEBUG] redisConfig : %v", redisConfig)
	client := redis.NewClient(&redis.Options{
		Addr: redisConfig.Addr,
		Password: redisConfig.Password,
		DB: redisConfig.Database,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
