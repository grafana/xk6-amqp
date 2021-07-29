package amqp

import (
	amqpDriver "github.com/streadway/amqp"
)

type Channels struct {
	Version    string
	Connection *amqpDriver.Connection
}

type ChannelOptions struct {
	ConnectionUrl string
}

type ChannelDeclareOptions struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqpDriver.Table
}

type ChannelBindOptions struct {
	DestinationChannelName string
	SourceChannelName      string
	RoutingKey              string
	NoWait                  bool
	Args                    amqpDriver.Table
}

type ChannelUnindOptions struct {
	DestinationChannelName string
	SourceChannelName      string
	RoutingKey              string
	NoWait                  bool
	Args                    amqpDriver.Table
}

func (Channels *Channels) Declare(options ChannelDeclareOptions) error {
	ch, err := Channels.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ChannelDeclare(
		options.Name,
		options.Kind,
		options.Durable,
		options.AutoDelete,
		options.Internal,
		options.NoWait,
		options.Args,
	)
}

func (Channels *Channels) Delete(name string) error {
	ch, err := Channels.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ChannelDelete(
		name,
		false, // ifUnused
		false, // noWait
	)
}

func (Channels *Channels) Bind(options ChannelBindOptions) error {
	ch, err := Channels.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ChannelBind(
		options.DestinationChannelName,
		options.RoutingKey,
		options.SourceChannelName,
		options.NoWait,
		options.Args,
	)
}

func (Channels *Channels) Unbind(options ChannelUnindOptions) error {
	ch, err := Channels.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ChannelUnbind(
		options.DestinationChannelName,
		options.RoutingKey,
		options.SourceChannelName,
		options.NoWait,
		options.Args,
	)
}
