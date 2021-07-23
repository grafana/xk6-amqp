package amqp

import (
	amqpDriver "github.com/streadway/amqp"
)

type Queue struct {
	Version    string
	Connection *amqpDriver.Connection
}

type QueueOptions struct {
	ConnectionUrl string
}

type DeclareOptions struct {
	Name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
	Args             amqpDriver.Table
}

type QueueBindOptions struct {
	QueueName    string
	ExchangeName string
	RoutingKey   string
	NoWait       bool
	Args         amqpDriver.Table
}

type QueueUnindOptions struct {
	QueueName    string
	ExchangeName string
	RoutingKey   string
	Args         amqpDriver.Table
}

func (queue *Queue) Declare(options DeclareOptions) (amqpDriver.Queue, error) {
	ch, err := queue.Connection.Channel()
	if err != nil {
		return amqpDriver.Queue{}, err
	}
	defer ch.Close()
	return ch.QueueDeclare(
		options.Name,
		options.Durable,
		options.DeleteWhenUnused,
		options.Exclusive,
		options.NoWait,
		options.Args,
	)
}

func (queue *Queue) Inspect(name string) (amqpDriver.Queue, error) {
	ch, err := queue.Connection.Channel()
	if err != nil {
		return amqpDriver.Queue{}, err
	}
	defer ch.Close()
	return ch.QueueInspect(name)
}

func (queue *Queue) Delete(name string) error {
	ch, err := queue.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDelete(
		name,
		false, // ifUnused
		false, // ifEmpty
		false, // noWait
	)
	return err
}

func (queue *Queue) Bind(options QueueBindOptions) error {
	ch, err := queue.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.QueueBind(
		options.QueueName,
		options.RoutingKey,
		options.ExchangeName,
		options.NoWait,
		options.Args,
	)
}

func (queue *Queue) Unbind(options QueueUnindOptions) error {
	ch, err := queue.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.QueueUnbind(
		options.QueueName,
		options.RoutingKey,
		options.ExchangeName,
		options.Args,
	)
}

func (queue *Queue) Purge(name string, noWait bool) (int, error) {
	ch, err := queue.Connection.Channel()
	if err != nil {
		return 0, err
	}
	defer ch.Close()
	return ch.QueuePurge(name, noWait)
}
