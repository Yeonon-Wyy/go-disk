package mq

import (
	"github.com/streadway/amqp"
	"go-disk/common/log4disk"
	"go-disk/services/transfer/config"
)

var (
	channel      *amqp.Channel
	consumerDone chan struct{}

	mqConfig = config.Conf.Mq
)

func initChannel() bool {
	if channel != nil {
		return true
	}

	conn, err := amqp.Dial(mqConfig.Rabbit.Url)
	if err != nil {
		log4disk.E("failed to connect to rabbit mq server : %v", err)
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		log4disk.E("failed to get rabbit mq channel : %v", err)
		return false
	}

	return true
}

func RabbitConsume(queueName string, consumerName string, callBack func([]byte) bool) {
	if !initChannel() {
		return
	}

	msgChannel, err := channel.Consume(
		queueName,
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log4disk.E("start consumer error : %v", err)
		return
	}

	consumerDone = make(chan struct{})

	go func() {
		for msg := range msgChannel {
			log4disk.D("consumer process success")
			if suc := callBack(msg.Body); !suc {
				//TODO: push another queue
			}
		}
	}()

	<-consumerDone
	//close rabbit channel
	channel.Close()
}
