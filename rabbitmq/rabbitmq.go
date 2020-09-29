package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"go-product/datamodels"
	"go-product/services"
	"log"
	"sync"
)

// rabbitmq connection information
const MQURL = "amqp://test:123456@192.168.56.38:5672/basic"

// rabbitMQ struct
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	// Queue name
	QueueName string
	// Exchange name
	Exchange string

	// bind key name
	Key string

	// Connection information
	Mqurl string
	sync.Mutex
}

// Create struct instance
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
}

// disconnection channel and connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// Error Handler
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// Create simple Rabbit instance
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	// create RabbitMQ instance
	rabbitmq := NewRabbitMQ(queueName, "", "")
	var err error

	// Get connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")

	// Get channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open channel")
	return rabbitmq
}

// simple producer
func (r *RabbitMQ) PublishSimple(message string) error {
	r.Lock()
	defer r.Unlock()
	// request queue if not exists create queue
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, //  persistence
		false, // auto delete
		false, // exclusive other
		false, // block
		nil,   // other arguments
	)

	if err != nil {
		return err
	}

	// call channel send message to queue
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false,
		false, // if set is true when not found consumer return message to sender
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return nil
}

// simple consumer
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
	}

	// Consume Qos limit
	r.channel.Qos(
		1,     // current consume once max message amount
		0,     // server delivery max size (8bit unit)
		false, // when set is true channel available
	)

	// Receive message
	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack; here we set is manual ack
		false,  // exclusive
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	// start coroutine
	go func() {
		for d := range msgs {
			// consume logic
			log.Printf("Received a message: %s", d.Body)
			message := &datamodels.Message{}
			// json decode
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}

			// Insert Order
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
			}

			// Reduce product number
			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				fmt.Println(err)
			}

			// When set true represent all not confirm message
			// When set false confirm current message
			d.Ack(false)
		}
	}()

	log.Printf("[*] Wating for messages. To exit press CTRL + C")
	<-forever

}

// Create pub/sub mode RabbitMQ instance
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	var err error

	// Get connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")

	// Get channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// pub/sub producer
func (r *RabbitMQ) PublishPub(message string) {
	// try to create exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false, // when internal is true exchange just using switch to other, client cannot use it
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	// send message
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// pub/sub consumer
func (r *RabbitMQ) RecieveSub() {
	// Try to create exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange!")

	// try to create queue
	q, err := r.channel.QueueDeclare(
		"", // random producer queue name
		false,
		false,
		true,
		false,
		nil,
	)

	// bind queue on exchange
	err = r.channel.QueueBind(
		q.Name,
		"", // pub/sub mode here kee is empty
		r.Exchange,
		false,
		nil,
	)

	// consume message
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("[*] Waiting for messages. To exit press CTRL + C")
	<-forever

}

// Create rabbitMQ routing mode instance
func NewRabbitMQRouting(exchangeName, routingKey string) *RabbitMQ {
	// create RabbitMQ instance
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error

	// Get connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")

	// Get channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open channel")
	return rabbitmq
}

// pub/sub producer
func (r *RabbitMQ) PublishRouting(message string) {
	// try to create exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange")

	// send message
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// Routing consumer
func (r *RabbitMQ) RecieveRouting() {
	// Try to create exchange
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an exchange!")

	// try to create queue
	q, err := r.channel.QueueDeclare(
		"", // random producer queue name
		false,
		false,
		true,
		false,
		nil,
	)

	// bind queue on exchange
	err = r.channel.QueueBind(
		q.Name,
		r.Key, // routing mode  key is required
		r.Exchange,
		false,
		nil,
	)

	// consume message
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("[*] Waiting for messages. To exit press CTRL + C")
	<-forever

}
