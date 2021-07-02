package amqp

import (
	amqpDriver "github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
)

const version = "v0.0.1"

type Amqp struct {
	Version    string
	Connection *amqpDriver.Connection
}

type AmqpOptions struct {
	ConnectionUrl string
}

type PublishOptions struct {
	QueueName string
	Body      string
}

func (amqp *Amqp) Start(options AmqpOptions) error {
	conn, err := amqpDriver.Dial(options.ConnectionUrl)
	amqp.Connection = conn
	return err
}

func (amqp *Amqp) Publish(options PublishOptions) error {
	ch, err := amqp.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	name := options.QueueName
	queue, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqpDriver.Publishing{
			ContentType: "text/plain",
			Body:        []byte(options.Body),
		},
	)
}

func init() {
	modules.Register("k6/x/amqp", &Amqp{
		Version: version,
	})
}
