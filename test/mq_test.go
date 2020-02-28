package test

import (
	"go-disk/config"
	"go-disk/mq"
	"testing"
)

func TestRabbitMQ(t *testing.T) {
	mq.RabbitPublish(config.RabbitExchangeName, config.RabbitCephRoutingKey, []byte("yeonon"))
}
