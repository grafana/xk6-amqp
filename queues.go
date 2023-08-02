package amqp

import (
	"fmt"

	amqpDriver "github.com/rabbitmq/amqp091-go"
)

// Queue defines a connection to a point-to-point destination.
type Queue struct {
	Version     string
	Connections *map[int]*amqpDriver.Connection
	MaxConnID   *int
}

// QueueOptions defines configuration settings for accessing a queue.
type QueueOptions struct {
	ConnectionURL string
}

// DeclareOptions provides queue options when declaring (creating) a queue.
type DeclareOptions struct {
	ConnectionID     int
	Name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
	Args             amqpDriver.Table
}

// QueueInspectOptions provide options when inspecting a queue.
type QueueInspectOptions struct {
	ConnectionID int
}

// QueueDeleteOptions provide options when deleting a queue.
type QueueDeleteOptions struct {
	ConnectionID int
}

// QueueBindOptions provides options when binding a queue to an exchange in order to receive message(s).
type QueueBindOptions struct {
	ConnectionID int
	QueueName    string
	ExchangeName string
	RoutingKey   string
	NoWait       bool
	Args         amqpDriver.Table
}

// QueueUnbindOptions provides options when unbinding a queue from an exchange to stop receiving message(s).
type QueueUnbindOptions struct {
	ConnectionID int
	QueueName    string
	ExchangeName string
	RoutingKey   string
	Args         amqpDriver.Table
}

// QueuePurgeOptions provide options when purging (emptying) a queue.
type QueuePurgeOptions struct {
	ConnectionID int
}

// GetConn gets an initialised connection by ID, or returns the last initialised one if ID is 0
func (queue *Queue) GetConn(connID int) (*amqpDriver.Connection, error) {
	if connID == 0 {
		conn := (*queue.Connections)[*queue.MaxConnID]
		if conn == nil {
			return &amqpDriver.Connection{}, fmt.Errorf("connection not initialised")
		}
		return conn, nil
	}

	conn := (*queue.Connections)[connID]
	if conn == nil {
		return &amqpDriver.Connection{}, fmt.Errorf("connection with ID %d not initialised", connID)
	}
	return conn, nil
}

// Declare creates a new queue given the provided options.
func (queue *Queue) Declare(options DeclareOptions) (amqpDriver.Queue, error) {
	conn, err := queue.GetConn(options.ConnectionID)
	if err != nil {
		return amqpDriver.Queue{}, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return amqpDriver.Queue{}, err
	}
	defer func() {
		_ = ch.Close()
	}()
	return ch.QueueDeclare(
		options.Name,
		options.Durable,
		options.DeleteWhenUnused,
		options.Exclusive,
		options.NoWait,
		options.Args,
	)
}

// Inspect provides queue metadata given queue name.
func (queue *Queue) Inspect(name string, options QueueInspectOptions) (amqpDriver.Queue, error) {
	conn, err := queue.GetConn(options.ConnectionID)
	if err != nil {
		return amqpDriver.Queue{}, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return amqpDriver.Queue{}, err
	}
	defer func() {
		_ = ch.Close()
	}()
	return ch.QueueInspect(name)
}

// Delete removes a queue from the remote server given the queue name.
func (queue *Queue) Delete(name string, options QueueDeleteOptions) error {
	conn, err := queue.GetConn(options.ConnectionID)
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
	_, err = ch.QueueDelete(
		name,
		false, // ifUnused
		false, // ifEmpty
		false, // noWait
	)
	return err
}

// Bind subscribes a queue to an exchange in order to receive message(s).
func (queue *Queue) Bind(options QueueBindOptions) error {
	conn, err := queue.GetConn(options.ConnectionID)
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
	return ch.QueueBind(
		options.QueueName,
		options.RoutingKey,
		options.ExchangeName,
		options.NoWait,
		options.Args,
	)
}

// Unbind removes a queue subscription from an exchange to discontinue receiving message(s).
func (queue *Queue) Unbind(options QueueUnbindOptions) error {
	conn, err := queue.GetConn(options.ConnectionID)
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
	return ch.QueueUnbind(
		options.QueueName,
		options.RoutingKey,
		options.ExchangeName,
		options.Args,
	)
}

// Purge removes all non-consumed message(s) from the specified queue.
func (queue *Queue) Purge(name string, noWait bool, options QueuePurgeOptions) (int, error) {
	conn, err := queue.GetConn(options.ConnectionID)
	if err != nil {
		return 0, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = ch.Close()
	}()
	return ch.QueuePurge(name, noWait)
}
