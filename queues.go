package amqp

import (
	amqpDriver "github.com/rabbitmq/amqp091-go"
)

// Queue defines a connection to a point-to-point destination.
type Queue struct {
	Version    string
	Connection *amqpDriver.Connection
}

// QueueOptions defines configuration settings for accessing a queue.
type QueueOptions struct {
	ConnectionURL string
}

// DeclareOptions provides queue options when declaring (creating) a queue.
type DeclareOptions struct {
	Name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
	Args             amqpDriver.Table
}

// QueueBindOptions provides options when binding a queue to an exchange in order to receive message(s).
type QueueBindOptions struct {
	QueueName    string
	ExchangeName string
	RoutingKey   string
	NoWait       bool
	Args         amqpDriver.Table
}

// QueueUnbindOptions provides options when unbinding a queue from an exchange to stop receiving message(s).
type QueueUnbindOptions struct {
	QueueName    string
	ExchangeName string
	RoutingKey   string
	Args         amqpDriver.Table
}

// Declare creates a new queue given the provided options.
func (queue *Queue) Declare(options DeclareOptions) (amqpDriver.Queue, error) {
	ch, err := queue.Connection.Channel()
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
func (queue *Queue) Inspect(name string) (amqpDriver.Queue, error) {
	ch, err := queue.Connection.Channel()
	if err != nil {
		return amqpDriver.Queue{}, err
	}
	defer func() {
		_ = ch.Close()
	}()
	return ch.QueueInspect(name)
}

// Delete removes a queue from the remote server given the queue name.
func (queue *Queue) Delete(name string) error {
	ch, err := queue.Connection.Channel()
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
	ch, err := queue.Connection.Channel()
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
	ch, err := queue.Connection.Channel()
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
func (queue *Queue) Purge(name string, noWait bool) (int, error) {
	ch, err := queue.Connection.Channel()
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = ch.Close()
	}()
	return ch.QueuePurge(name, noWait)
}
