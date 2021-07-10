package amqp

import (
	amqpDriver "github.com/streadway/amqp"
)

type Exchanges struct {
	Version    string
	Connection *amqpDriver.Connection
}

type ExchangeOptions struct {
	ConnectionUrl string
}

type EchangeDeclareOptions struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqpDriver.Table
}

func (exchanges *Exchanges) Declare(options EchangeDeclareOptions) error {
	ch, err := exchanges.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ExchangeDeclare(
		options.Name,
		options.Kind,
		options.Durable,
		options.AutoDelete,
		options.Internal,
		options.NoWait,
		options.Args,
	)
}

func (exchanges *Exchanges) Delete(name string) error {
	ch, err := exchanges.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ExchangeDelete(
		name,
		false, // ifUnused
		false, // noWait
	)
}
