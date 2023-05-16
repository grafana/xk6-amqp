// Package amqp contains AMQP API for a remote server.
package amqp

import (
	"context"
	"encoding/json"
	"time"
	"fmt"

	amqpDriver "github.com/rabbitmq/amqp091-go"
	"github.com/vmihailenco/msgpack/v5"
	"go.k6.io/k6/js/modules"
)

const version = "v0.3.0"

// AMQP type holds connection to a remote AMQP server.
type AMQP struct {
	Version     string
	Connections *map[int]*amqpDriver.Connection
	MaxConnId   *int
	Queue       *Queue
	Exchange    *Exchange
}

// Options defines configuration options for an AMQP session.
type Options struct {
	ConnectionURL string
}

// PublishOptions defines a message payload with delivery options.
type PublishOptions struct {
	ConnectionId  int
	QueueName     string
	Body          string
	Headers       amqpDriver.Table
	Exchange      string
	ContentType   string
	Mandatory     bool
	Immediate     bool
	Persistent    bool
	CorrelationId string
	ReplyTo       string
	Expiration    string
	MessageId     string
	Timestamp     int64 // unix epoch timestamp in seconds
	Type          string
	UserId        string
	AppId         string
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
	ConnectionId int
	Listener     ListenerType
	QueueName    string
	Consumer     string
	AutoAck      bool
	Exclusive    bool
	NoLocal      bool
	NoWait       bool
	Args         amqpDriver.Table
}

const messagepack = "application/x-msgpack"

// Start establishes a session with an AMQP server given the provided options.
func (amqp *AMQP) Start(options Options) (int, error) {
	conn, err := amqpDriver.Dial(options.ConnectionURL)
	*amqp.MaxConnId += 1
	(*amqp.Connections)[*amqp.MaxConnId] = conn
	return *amqp.MaxConnId, err
}

// Gets an initialised connection by ID, or returns the last initialised one if ID is 0
func (amqp *AMQP) GetConn(connId int) (*amqpDriver.Connection, error) {
	if connId == 0 {
		conn := (*amqp.Connections)[*amqp.MaxConnId]
		if conn == nil {
			return &amqpDriver.Connection{}, fmt.Errorf("Connection not initialised")
		}
		return conn, nil
	} else {
		conn := (*amqp.Connections)[connId]
		if conn == nil {
			return &amqpDriver.Connection{}, fmt.Errorf("Connection with ID %d not initialised", connId)
		}
		return conn, nil
	}
}

// Publish delivers the payload using options provided.
func (amqp *AMQP) Publish(options PublishOptions) error {
	conn, err := amqp.GetConn(options.ConnectionId)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()

	publishing := amqpDriver.Publishing{
		Headers:     options.Headers,
		ContentType: options.ContentType,
	}

	if options.ContentType == messagepack {
		var jsonParsedBody interface{}

		if err = json.Unmarshal([]byte(options.Body), &jsonParsedBody); err != nil {
			return err
		}

		publishing.Body, err = msgpack.Marshal(jsonParsedBody)
		if err != nil {
			return err
		}
	} else {
		publishing.Body = []byte(options.Body)
	}

	if options.Persistent {
		publishing.DeliveryMode = amqpDriver.Persistent
	}

	// Well, I guess 1970-01-01 isn't allowed now.
	if options.Timestamp != 0 {
		publishing.Timestamp = time.Unix(options.Timestamp, 0)
	}

	publishing.CorrelationId = options.CorrelationId
	publishing.ReplyTo = options.ReplyTo
	publishing.Expiration = options.Expiration
	publishing.MessageId = options.MessageId
	publishing.Type = options.Type
	publishing.UserId = options.UserId
	publishing.AppId = options.AppId

	return ch.PublishWithContext(
		context.Background(), // TODO: use vu context
		options.Exchange,
		options.QueueName,
		options.Mandatory,
		options.Immediate,
		publishing,
	)
}

// Listen binds to an AMQP queue in order to receive message(s) as they are received.
func (amqp *AMQP) Listen(options ListenOptions) error {
	conn, err := amqp.GetConn(options.ConnectionId)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
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
	connections := make(map[int]*amqpDriver.Connection)
	maxConnId := 0
	queue := Queue{
		Connections: &connections,
		MaxConnId:   &maxConnId,
	}
	exchange := Exchange{
		Connections: &connections,
		MaxConnId:   &maxConnId,
	}
	generalAMQP := AMQP{
		Version:     version,
		Connections: &connections,
		MaxConnId:   &maxConnId,
		Queue:       &queue,
		Exchange:    &exchange,
	}

	modules.Register("k6/x/amqp", &generalAMQP)
	modules.Register("k6/x/amqp/queue", &queue)
	modules.Register("k6/x/amqp/exchange", &exchange)
}
