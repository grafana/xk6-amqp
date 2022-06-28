package amqp

import (
	amqpDriver "github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
)

const version = "v0.0.1"

type Amqp struct {
	Version    string
	Connection *amqpDriver.Connection
	Queue     *Queue
	Exchange  *Exchange
}

type AmqpOptions struct {
	ConnectionUrl string
}

type PublishOptions struct {
	QueueName string
	Body      string
	Exchange  string
	ContentType string
	Mandatory bool
	Immediate bool
}

type ConsumeOptions struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqpDriver.Table
}

type ListenerType func(string) error

type ListenOptions struct {
	Listener  ListenerType
	QueueName string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqpDriver.Table
}

func (amqp *Amqp) Start(options AmqpOptions) error {
	conn, err := amqpDriver.Dial(options.ConnectionUrl)
	amqp.Connection = conn
	amqp.Queue.Connection = conn
	amqp.Exchange.Connection = conn
	return err
}

func (amqp *Amqp) Publish(options PublishOptions) error {
	ch, err := amqp.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(
		options.Exchange,
		options.QueueName,
		options.Mandatory,
		options.Immediate,
		amqpDriver.Publishing{
			ContentType: options.ContentType,
			Body:        []byte(options.Body),
		},
	)
}

func (amqp *Amqp) Listen(options ListenOptions) error {
	ch, err := amqp.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		options.QueueName,
		options.Consumer,
		options.AutoAck,
		options.Exclusive,
		options.NoLocal,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			options.Listener(string(d.Body))
		}
	}()
	return nil
}

func init() {

	queue := Queue{}
	exchange := Exchange{}
	generalAmqp := Amqp{
		Version:   version,
		Queue:    &queue,
		Exchange: &exchange,
	}

	modules.Register("k6/x/amqp", &generalAmqp)
	modules.Register("k6/x/amqp/queue", &queue)
	modules.Register("k6/x/amqp/exchange", &exchange)
}
