package mq

import (
	"github.com/streadway/amqp"
	"go-disk/config"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel

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

func Publish(exchange string, routingKey string, msg []byte) bool {
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

func Consume() {

}