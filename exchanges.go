package amqp

import (
	"fmt"
	amqpDriver "github.com/rabbitmq/amqp091-go"
)

// Exchange defines a connection to publish/subscribe destinations.
type Exchange struct {
	Version     string
	Connections *map[int]*amqpDriver.Connection
	MaxConnId   *int
}

// ExchangeOptions defines configuration settings for accessing an exchange.
type ExchangeOptions struct {
	ConnectionURL string
}

// ExchangeDeclareOptions provides options when declaring (creating) an exchange.
type ExchangeDeclareOptions struct {
	ConnectionId int
	Name         string
	Kind         string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	NoWait       bool
	Args         amqpDriver.Table
}

// ExchangeDeleteOptions provides options when deleting an exchange.
type ExchangeDeleteOptions struct {
	ConnectionId int
}

// ExchangeBindOptions provides options when binding (subscribing) one exchange to another.
type ExchangeBindOptions struct {
	ConnectionId            int
	DestinationExchangeName string
	SourceExchangeName      string
	RoutingKey              string
	NoWait                  bool
	Args                    amqpDriver.Table
}

// ExchangeUnbindOptions provides options when unbinding (unsubscribing) one exchange from another.
type ExchangeUnbindOptions struct {
	ConnectionId            int
	DestinationExchangeName string
	SourceExchangeName      string
	RoutingKey              string
	NoWait                  bool
	Args                    amqpDriver.Table
}

// Gets an initialised connection by ID, or returns the last initialised one if ID is 0
func (exchange *Exchange) GetConn(connId int) (*amqpDriver.Connection, error) {
	if connId == 0 {
		conn := (*exchange.Connections)[*exchange.MaxConnId]
		if conn == nil {
			return &amqpDriver.Connection{}, fmt.Errorf("Connection not initialised")
		}
		return conn, nil
	} else {
		conn := (*exchange.Connections)[connId]
		if conn == nil {
			return &amqpDriver.Connection{}, fmt.Errorf("Connection with ID %d not initialised", connId)
		}
		return conn, nil
	}
}

// Declare creates a new exchange given the provided options.
func (exchange *Exchange) Declare(options ExchangeDeclareOptions) error {
	conn, err := exchange.GetConn(options.ConnectionId)
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
func (exchange *Exchange) Delete(name string, options ExchangeDeleteOptions) error {
	conn, err := exchange.GetConn(options.ConnectionId)
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
	return ch.ExchangeDelete(
		name,
		false, // ifUnused
		false, // noWait
	)
}

// Bind subscribes one exchange to another.
func (exchange *Exchange) Bind(options ExchangeBindOptions) error {
	conn, err := exchange.GetConn(options.ConnectionId)
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
	conn, err := exchange.GetConn(options.ConnectionId)
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
	return ch.ExchangeUnbind(
		options.DestinationExchangeName,
		options.RoutingKey,
		options.SourceExchangeName,
		options.NoWait,
		options.Args,
	)
}
