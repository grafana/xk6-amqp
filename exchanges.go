package amqp

import (
	amqpDriver "github.com/streadway/amqp"
)

// Exchange defines a connection to publish/subscribe destinations.
type Exchange struct {
	Version    string
	Connection *amqpDriver.Connection
}

// ExchangeOptions defines configuration settings for accessing an exchange.
type ExchangeOptions struct {
	ConnectionURL string
}

// ExchangeDeclareOptions provides options when declaring (creating) an exchange.
type ExchangeDeclareOptions struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqpDriver.Table
}

// ExchangeBindOptions provides options when binding (subscribing) one exchange to another.
type ExchangeBindOptions struct {
	DestinationExchangeName string
	SourceExchangeName      string
	RoutingKey              string
	NoWait                  bool
	Args                    amqpDriver.Table
}

// ExchangeUnbindOptions provides options when unbinding (unsubscribing) one exchange from another.
type ExchangeUnbindOptions struct {
	DestinationExchangeName string
	SourceExchangeName      string
	RoutingKey              string
	NoWait                  bool
	Args                    amqpDriver.Table
}

// Declare creates a new exchange given the provided options.
func (exchange *Exchange) Declare(options ExchangeDeclareOptions) error {
	ch, err := exchange.Connection.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
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

// Delete removes an exchange from the remote server given the exchange name.
func (exchange *Exchange) Delete(name string) error {
	ch, err := exchange.Connection.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
	return ch.ExchangeDelete(
		name,
		false, // ifUnused
		false, // noWait
	)
}

// Bind subscribes one exchange to another.
func (exchange *Exchange) Bind(options ExchangeBindOptions) error {
	ch, err := exchange.Connection.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
	return ch.ExchangeBind(
		options.DestinationExchangeName,
		options.RoutingKey,
		options.SourceExchangeName,
		options.NoWait,
		options.Args,
	)
}

// Unbind removes a subscription from one exchange to another.
func (exchange *Exchange) Unbind(options ExchangeUnbindOptions) error {
	ch, err := exchange.Connection.Channel()
	if err != nil {
		return err
	}
	defer func() {
		_ = ch.Close()
	}()
	return ch.ExchangeUnbind(
		options.DestinationExchangeName,
		options.RoutingKey,
		options.SourceExchangeName,
		options.NoWait,
		options.Args,
	)
}
