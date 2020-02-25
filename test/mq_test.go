package test

import (
	"go-disk/config"
	"go-disk/mq"
	"testing"
)

func TestRabbitMQ(t *testing.T) {
	mq.Publish(config.RabbitExchangeName, config.RabbitCephRoutingKey, []byte("yeonon"))
}
