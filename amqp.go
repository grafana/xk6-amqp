// Package amqp contains AMQP API for a remote server.
package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqpDriver "github.com/rabbitmq/amqp091-go"
	"github.com/vmihailenco/msgpack/v5"
	"go.k6.io/k6/js/modules"
)

const version = "v0.4.1"

// AMQP type holds connection to a remote AMQP server.
type AMQP struct {
	Version     string
	Connections *map[int]*amqpDriver.Connection
	MaxConnID   *int
	Queue       *Queue
	Exchange    *Exchange
}

// Options defines configuration options for an AMQP session.
type Options struct {
	ConnectionURL string
}

// PublishOptions defines a message payload with delivery options.
type PublishOptions struct {
	ConnectionID  int
	QueueName     string
	Body          string
	Headers       amqpDriver.Table
	Exchange      string
	ContentType   string
	Mandatory     bool
	Immediate     bool
	Persistent    bool
	CorrelationID string
	ReplyTo       string
	Expiration    string
	MessageID     string
	Timestamp     int64 // unix epoch timestamp in seconds
	Type          string
	UserID        string
	AppID         string
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
	ConnectionID int
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
	*amqp.MaxConnID++
	(*amqp.Connections)[*amqp.MaxConnID] = conn
	return *amqp.MaxConnID, err
}

// GetConn gets an initialised connection by ID, or returns the last initialised one if ID is 0
func (amqp *AMQP) GetConn(connID int) (*amqpDriver.Connection, error) {
	if connID == 0 {
		conn := (*amqp.Connections)[*amqp.MaxConnID]
		if conn == nil {
			return &amqpDriver.Connection{}, fmt.Errorf("connection not initialised")
		}
		return conn, nil
	}

	conn := (*amqp.Connections)[connID]
	if conn == nil {
		return &amqpDriver.Connection{}, fmt.Errorf("connection with ID %d not initialised", connID)
	}
	return conn, nil
}

// Publish delivers the payload using options provided.
func (amqp *AMQP) Publish(options PublishOptions) error {
	conn, err := amqp.GetConn(options.ConnectionID)
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

	publishing.CorrelationId = options.CorrelationID
	publishing.ReplyTo = options.ReplyTo
	publishing.Expiration = options.Expiration
	publishing.MessageId = options.MessageID
	publishing.Type = options.Type
	publishing.UserId = options.UserID
	publishing.AppId = options.AppID

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
	conn, err := amqp.GetConn(options.ConnectionID)
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
	maxConnID := 0
	queue := Queue{
		Connections: &connections,
		MaxConnID:   &maxConnID,
	}
	exchange := Exchange{
		Connections: &connections,
		MaxConnID:   &maxConnID,
	}
	generalAMQP := AMQP{
		Version:     version,
		Connections: &connections,
		MaxConnID:   &maxConnID,
		Queue:       &queue,
		Exchange:    &exchange,
	}

	modules.Register("k6/x/amqp", &generalAMQP)
	modules.Register("k6/x/amqp/queue", &queue)
	modules.Register("k6/x/amqp/exchange", &exchange)
}
