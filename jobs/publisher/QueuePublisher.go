package publisher

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gostk/logger"
	"strconv"
)

func Publish(input interface{}, queue string) bool {
	//init rabbitmq
	status := true
	logger.Log.Debugw("Initialising RabbitMQ Publish ", "queue", queue, "payload", input)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		status = false
		logger.Log.Errorw("RabbitMQ Publish Failed Initializing Broker Connection ", "error", err, "queue", queue)
		panic(err)
	}

	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	ch, err := conn.Channel()
	if err != nil {
		status = false
		logger.Log.Errorw("RabbitMQ Publish Connection Error ", "error", err, "queue", queue)
	}
	defer ch.Close()

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	// We can print out the status of our Queue here
	// this will information like the amount of messages on
	// the queue
	logger.Log.Debugw("RabbitMQ Publish Queue Status ", "status", q, "queue", queue, "payload", input)
	// Handle any errors if we were unable to create the queue
	if err != nil {
		status = false
		fmt.Println(err)
	}

	inputBytes, err := json.Marshal(input)
	if err != nil {
		status = false
		logger.Log.Errorw("RabbitMQ Publish Payload Marshal Error ", "error", err, "queue", queue)
		return status
	}

	// attempt to publish a message to the queue!
	err = ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        inputBytes,
		},
	)

	if err != nil {
		status = false
		logger.Log.Errorw("RabbitMQ Publish Payload Publish Error ", "error", err, "queue", queue)
		fmt.Println(err)
	}

	logger.Log.Debugw("RabbitMQ Publish Successfully Published Message to Queue ~> "+strconv.FormatBool(status), "queue", queue, "payload", input)
	return status
}
