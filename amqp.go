// Package amqp contains AMQP API for a remote server.
package amqp

import (
	amqpDriver "github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

const version = "v0.1.0"

// AMQP type holds connection to a remote AMQP server.
type AMQP struct {
	Version    string
	Connection *amqpDriver.Connection
	Queue      *Queue
	Exchange   *Exchange
}

// Options defines configuration options for an AMQP session.
type Options struct {
	ConnectionURL string
}

// PublishOptions defines a message payload with delivery options.
type PublishOptions struct {
	QueueName   string
	Body        string
	Exchange    string
	ContentType string
	Mandatory   bool
	Immediate   bool
	Persistent  bool
	CypherMethod string
}

// ConsumeOptions defines options for use when consuming a message.
type ConsumeOptions struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqpDriver.Table
}

// ListenerType is the message handler implemented within JavaScript.
type ListenerType func(string) error

// ListenOptions defines options for subscribing to message(s) within a queue.
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

// Start establishes a session with an AMQP server given the provided options.
func (amqp *AMQP) Start(options Options) error {
	conn, err := amqpDriver.Dial(options.ConnectionURL)
	amqp.Connection = conn
	amqp.Queue.Connection = conn
	amqp.Exchange.Connection = conn
	return err
}

// Publish delivers the payload using options provided.
func (amqp *AMQP) Publish(options PublishOptions) error {
	ch, err := amqp.Connection.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()

  if options.CypherMethod == "messagepack" {
    cypherBody, err := msgpack.Marshal(options.Body)
  }

  if err != nil {
    panic(err)
  }

	publishing := amqpDriver.Publishing{
		ContentType: options.ContentType,
		Body:        []byte(cypherBody),
	}

	if options.Persistent {
		publishing.DeliveryMode = amqpDriver.Persistent
	}

	return ch.Publish(
		options.Exchange,
		options.QueueName,
		options.Mandatory,
		options.Immediate,
		publishing,
	)
}

// Listen binds to an AMQP queue in order to receive message(s) as they are received.
func (amqp *AMQP) Listen(options ListenOptions) error {
	ch, err := amqp.Connection.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()

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
			err = options.Listener(string(d.Body))
		}
	}()
	return err
}

func init() {
	queue := Queue{}
	exchange := Exchange{}
	generalAMQP := AMQP{
		Version:  version,
		Queue:    &queue,
		Exchange: &exchange,
	}

	modules.Register("k6/x/amqp", &generalAMQP)
	modules.Register("k6/x/amqp/queue", &queue)
	modules.Register("k6/x/amqp/exchange", &exchange)
}
