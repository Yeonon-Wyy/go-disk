package mq

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"go-disk/config"
	"log"
)

var channel *amqp.Channel
var consumerDone chan struct{}


func initChannel() bool {
	if channel != nil {
		return true
	}

	conn, err := amqp.Dial(config.RabbitUrl)
	if err != nil {
		log.Printf("failed to connect to rabbit mq server : %v", err)
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Printf("failed to get rabbit mq channel : %v", err)
		return false
	}

	return true
}

func RabbitPublish(exchange string, routingKey string, msg []byte) bool {
	if !initChannel() {
		return false
	}

	err := channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: msg,
		})

	if err != nil {
		log.Printf("publish message error : %v", err)
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
		log.Printf("start consumer error : %v", err)
		return
	}

	consumerDone = make(chan struct{})

	go func() {
		for msg := range msgChannel {
			log.Println("consumer process success")
			if suc := callBack(msg.Body); !suc {
				//TODO: push another queue
			}
		}
	}()

	<-consumerDone
	//close rabbit channel
	channel.Close()
}


func PublishError(exchange string, routingKey string, errMsg RabbitErrMessage) error {
	data, err := json.Marshal(errMsg)
	if err != nil {
		log.Printf("json marshal error : %v", err)
		return err
	}
	if suc := RabbitPublish(exchange, routingKey, data); !suc {
		return errors.New("failed to publish message to rabbit mq ")
	}
	return nil
}