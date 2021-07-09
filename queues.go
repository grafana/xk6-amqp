package amqp

import (
	amqpDriver "github.com/streadway/amqp"
)

type Queues struct {
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

func (queues *Queues) Declare(options DeclareOptions) (amqpDriver.Queue, error) {
	ch, err := queues.Connection.Channel()
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

func (queues *Queues) Inspect(name string) (amqpDriver.Queue, error) {
	ch, err := queues.Connection.Channel()
	if err != nil {
		return amqpDriver.Queue{}, err
	}
	defer ch.Close()
	return ch.QueueInspect(name)
}

func (queues *Queues) Delete(name string) error {
	ch, err := queues.Connection.Channel()
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
